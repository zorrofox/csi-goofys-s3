package s3

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "io"
    "fmt"
    "os"

    "github.com/golang/glog"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"

    s3Native "github.com/aws/aws-sdk-go/service/s3"
)


const (
    metadataName = ".metadata.json"
    fsPrefix     = "csi-fs"
)

type s3Client struct {
    cfg   *Config
    native *s3Native.S3
}

type bucket struct {
    Name          string
    Mounter       string
    FSPath        string
    CapacityBytes int64
}

func newS3Client(cfg *Config) (*s3Client, error) {
    var client = &s3Client{}

    client.cfg = cfg

    if cfg.Region == "" {
        return nil, fmt.Errorf("Failed to create S3 Client for no region!")
    }

    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(cfg.Region),
        },
    )

    nativeClient := s3Native.New(sess)

    if err != nil {
        return nil, err
    }
    client.native = nativeClient
    return client, nil
}


func newS3NativeClient() (*s3Client, error) {
    var region string
    if os.Getenv("AWS_REGION") !=""{

        region = os.Getenv("AWS_REGION")
    } else if os.Getenv("AWS_DEFAULT_REGION") !=""{
        region = os.Getenv("AWS_DEFAULT_REGION")
    } else {
        region = defaultRegion
    }
    glog.V(3).Infof("S3 Client Using the Region: %s", region)
    return newS3Client(&Config{
        Region: region,
        Mounter: "",
        ReadOnly: false,
    })
}

func (client *s3Client) bucketExists(bucketName string) (bool, error) {
    _, err := client.native.HeadBucket(&s3Native.HeadBucketInput{
        Bucket: aws.String(bucketName), // Required
    })

    if err != nil {
        glog.V(3).Infof("S3 Bucket Head Error: %v", err)
        return false, nil
    } else{
        return true, nil
    }

}

func (client *s3Client) createBucket(bucketName string) error {
    _, err := client.native.CreateBucket(&s3Native.CreateBucketInput{
        Bucket: aws.String(bucketName),
    })
    return err
}

func (client *s3Client) createPrefix(bucketName string, prefix string) error {
    _, err := client.native.PutObject(&s3Native.PutObjectInput{
        Bucket: aws.String(bucketName),
        Key: aws.String(prefix+"/"),
        Body: aws.ReadSeekCloser(bytes.NewReader([]byte(""))),
    })

    return err
}

func (client *s3Client) removeBucket(bucketName string) error {
    err := client.emptyBucket(bucketName)
    if err != nil {
        return err
    }
    _, err = client.native.DeleteBucket(&s3Native.DeleteBucketInput{
        Bucket: aws.String(bucketName),
    })
    return err
}

func (client *s3Client) emptyBucket(bucketName string) error {

    iter := s3manager.NewDeleteListIterator(client.native, &s3Native.ListObjectsInput{
        Bucket: aws.String(bucketName),
    })

    if err := s3manager.NewBatchDeleteWithClient(client.native).Delete(aws.BackgroundContext(), iter); err != nil {
        glog.Errorf("Unable to delete objects from bucket %s, error: %s", bucketName, err)
        return err
    }

    return nil
}

func (client *s3Client) setBucket(bucket *bucket) error {
    b := new(bytes.Buffer)
    json.NewEncoder(b).Encode(bucket)

    _, err := client.native.PutObject(&s3Native.PutObjectInput{
        Bucket: aws.String(bucket.Name),
        Key: aws.String(metadataName),
        Body: bytes.NewReader(b.Bytes()),
        ContentLength: aws.Int64(int64(b.Len())),
    })

    return err
}

func (client *s3Client) getBucket(bucketName string) (*bucket, error) {

    obj, err := client.native.GetObject(&s3Native.GetObjectInput{
        Bucket: aws.String(bucketName),
        Key: aws.String(metadataName),
    })

    if err != nil {
        return &bucket{}, err
    }

    b, err := ioutil.ReadAll(obj.Body)

    if err != nil && err != io.EOF {
        return &bucket{}, err
    }
    var meta bucket
    err = json.Unmarshal(b, &meta)
    return &meta, err
}

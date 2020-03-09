package s3

// Config holds values to configure the driver
type Config struct {
	Region   string
	Mounter  string
	ReadOnly bool
}

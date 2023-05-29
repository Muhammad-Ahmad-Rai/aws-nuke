package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	log "github.com/sirupsen/logrus"
)

type SESReceiptFilter struct {
	svc  *ses.SES
	name *string
}

func init() {
	register("SESReceiptFilter", ListSESReceiptFilters)
}

func ListSESReceiptFilters(sess *session.Session) ([]Resource, error) {
	svc := ses.New(sess)
	resources := []Resource{}

	params := &ses.ListReceiptFiltersInput{}

	output, err := svc.ListReceiptFilters(params)
	time.Sleep(1 * time.Second)
	
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "InvalidAction" {
			log.Debugf("skipping list operation for SESReceiptFilter: %s", awsErr.Message())
			// AWS responds with InvalidAction on regions that do not
			// support ListSESReceiptFilters.
			return resources,nil
		}

		return nil, err
	}

	for _, filter := range output.Filters {
		resources = append(resources, &SESReceiptFilter{
			svc:  svc,
			name: filter.Name,
		})
	}

	return resources, nil
}

func (f *SESReceiptFilter) Remove() error {

	_, err := f.svc.DeleteReceiptFilter(&ses.DeleteReceiptFilterInput{
		FilterName: f.name,
	})

	return err
}

func (f *SESReceiptFilter) String() string {
	return *f.name
}

package collector

import (
	"encoding/xml"
	"fmt"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var pmDesc = prometheus.NewDesc(
	prometheus.BuildFQName(namespace, "pm", "data"),
	"THIS IS DESCRIPTION TEXT FOR PROMETHEUS METRICS",
	[]string{
		"SAMPLETEST1",
	},
	nil,
)

type pmCollector struct {
	logger   log.Logger
}
type measDataFile struct {
	XMLName    xml.Name   `xml:"measDataFile"`
	FileHeader fileHeader `xml:"fileHeader"`
	MeasData   measData   `xml:"measData"`
	FileFooter fileFooter `xml:"fileFooter"`
}
type fileHeader struct {
	XMLName           xml.Name   `xml:"fileHeader"`
	FileFormatVersion string     `xml:"fileFormatVersion,attr"`
	VendorName        string     `xml:"vendorName,attr"`
	DnPrefix          string     `xml:"dnPrefix,attr"`
	FileSender        fileSender `xml:"fileSender"`
	MeasData          measData   `xml:"measData"`
}
type fileSender struct {
	XMLName    xml.Name `xml:"fileSender"`
	SenderName string   `xml:"senderName,attr"`
	SenderType string   `xml:"senderType,attr"`
}
type fileFooter struct {
	XMLName  xml.Name `xml:"fileFooter"`
	MeasData measData `xml:"measData"`
}
type measData struct {
	XMLName    xml.Name   `xml:"measData"`
	MeasEntity measEntity `xml:"measEntity"`
	MeasInfo   []measInfo `xml:"measInfo"`
	BeginTime  string     `xml:"beginTime,attr"`
	EndTime    string     `xml:"endTime,attr"`
}
type measEntity struct {
	XMLName xml.Name `xml:"measEntity"`
	Key     string   `xml:"localDn,attr"`
	Key2    string   `xml:"swVersion,attr"`
}
type measInfo struct {
	XMLName    xml.Name   `xml:"measInfo"`
	MeasInfoID string     `xml:"measInfoId,attr"`
	Job        job        `xml:"job"`
	GranPeriod granPeriod `xml:"granPeriod"`
	RepPeriod  repPeriod  `xml:"repPeriod"`
	MeasType   []measType `xml:"measType"`
	MeasValue  measValue  `xml:"measValue"`
}
type job struct {
	XMLName xml.Name `xml:"job"`
	//XMLAttr xml.Attr `xml:"jobId,attr"`
	Key string `xml:"jobId,attr"`
}
type granPeriod struct {
	XMLName xml.Name `xml:"granPeriod"`
	Key     string   `xml:"duration,attr"`
	Key2    string   `xml:"endTime,attr"`
}
type repPeriod struct {
	XMLName xml.Name `xml:"repPeriod"`
	Key     string   `xml:"duration,attr"`
}
type measType struct {
	XMLName xml.Name `xml:"measType"`
	Key     string   `xml:"p,attr"`
	Value   string   `xml:",chardata"`
}
type measValue struct {
	XMLName xml.Name `xml:"measValue"`
	Key     string   `xml:"measObjLdn,attr"`
	R       []r      `xml:"r"`
}
type r struct {
	XMLName xml.Name `xml:"r"`
	Key     string   `xml:"p,attr"`
	Value   float64   `xml:",chardata"`
}


func init() {
	registerCollector("pm", defaultEnabled, newPmCollector)
}

// newPmCollector returns new pmCollector.
func newPmCollector(logger log.Logger) (Collector, error) {
	return &pmCollector{logger}, nil
}
var fileNumber = 0

func (c *pmCollector) Update(ch chan<- prometheus.Metric) error {
	fmt.Println("update...")
	// xml 파일 오픈
	fileNumber += 1
	if fileNumber>5 {
		fileNumber=1
	}
	fp, err := os.Open("/go/parse_this_"+strconv.Itoa(fileNumber)+".xml")
	fmt.Println("Opening file : "+strconv.Itoa(fileNumber))
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	// xml 파일 읽기
	data, err := ioutil.ReadAll(fp)
	
	// xml 디코딩
	var measDataFile measDataFile
	err = xml.Unmarshal(data, &measDataFile)
	if err != nil {
		panic(err)
	}

	ch <- prometheus.MustNewConstMetric(pmDesc, prometheus.GaugeValue, measDataFile.MeasData.MeasInfo[0].MeasValue.R[0].Value,
		measDataFile.FileFooter.XMLName.Local,
	)

	measInfoList := measDataFile.MeasData.MeasInfo
	measInfoListLen := len(measInfoList)

	for i := 0; i < measInfoListLen; i++ {
		measTypeList := measDataFile.MeasData.MeasInfo[i].MeasType
		//measInfoIdValue := measDataFile.MeasData.MeasInfo[i].MeasInfoID
		measTypeListLen := len(measTypeList)
		//fmt.Println(measInfoListLen, measTypeListLen)

		for j := 0; j < measTypeListLen; j++ {
			metricKey := strings.ToLower(strings.ReplaceAll(measTypeList[j].Value, ".", "_"))
			metricValue := measInfoList[i].MeasValue.R[j].Value
			ch <- prometheus.MustNewConstMetric(
				pmDesc,
				prometheus.GaugeValue,
				metricValue,
				metricKey,
			)
		}
	}
	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

//Usage is the outside object from Twilios Usage api
type Usage struct {
	FirstPageURI    string         `json:"first_page_uri"`
	End             int            `json:"end"`
	PreviousPageURI interface{}    `json:"previous_page_uri"`
	URI             string         `json:"uri"`
	PageSize        int            `json:"page_size"`
	Start           int            `json:"start"`
	UsageRecords    []UsageRecords `json:"usage_records"`
	NextPageURI     string         `json:"next_page_uri"`
	Page            int            `json:"page"`
}

//UsageRecords are the inside objects from Twilios Usage api
type UsageRecords struct {
	Category    string    `json:"category"`
	Description string    `json:"description"`
	AccountSid  string    `json:"account_sid"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
	AsOf        time.Time `json:"as_of"`
	Count       float64   `json:",string"`
	CountUnit   string    `json:"count_unit"`
	Usage       float64   `json:",string"`
	UsageUnit   string    `json:"usage_unit"`
	Price       float64   `json:",string"`
	PriceUnit   string    `json:"price_unit"`
	APIVersion  string    `json:"api_version"`
	URI         string    `json:"uri"`
}

//UsageCollector creates the base Description objects for Prometheus Metrics
type UsageCollector struct {
	callerIDLookups         *prometheus.Desc
	calls                   *prometheus.Desc
	callsClient             *prometheus.Desc
	callsSip                *prometheus.Desc
	callsInbound            *prometheus.Desc
	callsInboundLocal       *prometheus.Desc
	callsInboundMobile      *prometheus.Desc
	callsInboundTollFree    *prometheus.Desc
	callsOutbound           *prometheus.Desc
	phoneNumbers            *prometheus.Desc
	phoneNumbersMobile      *prometheus.Desc
	phoneNumbersLocal       *prometheus.Desc
	phoneNumbersTollFree    *prometheus.Desc
	shortCodes              *prometheus.Desc
	shortCodesCustomerOwned *prometheus.Desc
	shortCodesRandom        *prometheus.Desc
	shortCodesVanity        *prometheus.Desc
	sms                     *prometheus.Desc
	smsInbound              *prometheus.Desc
	smsInboundLongCode      *prometheus.Desc
	smsInboundShortCode     *prometheus.Desc
	smsOutbound             *prometheus.Desc
	smsOutboundLongCode     *prometheus.Desc
	smsOutboundShortCode    *prometheus.Desc
	mms                     *prometheus.Desc
	mmsInbound              *prometheus.Desc
	mmsInboundLongCode      *prometheus.Desc
	mmsInboundShortCode     *prometheus.Desc
	mmsOutbound             *prometheus.Desc
	mmsOutboundLongCode     *prometheus.Desc
	mmsOutboundShortCode    *prometheus.Desc
	recordings              *prometheus.Desc
	recordingsStorage       *prometheus.Desc
	transcriptions          *prometheus.Desc
	mediaStorage            *prometheus.Desc
	authySMSOutbound        *prometheus.Desc
	authyCallsOutbound      *prometheus.Desc
	authyAuthentications    *prometheus.Desc
	authyPhoneVerifications *prometheus.Desc
	authyPhoneIntelligence  *prometheus.Desc
	authyMonthlyFees        *prometheus.Desc
	monitorStorage          *prometheus.Desc
	monitorReads            *prometheus.Desc
	monitorWrites           *prometheus.Desc
	taskRouterTasks         *prometheus.Desc
	turnMegabytes           *prometheus.Desc
	callRecordings          *prometheus.Desc
	trunkingRecordings      *prometheus.Desc
	trunkingTermination     *prometheus.Desc
	trunkingOrigination     *prometheus.Desc
}

//newUsageCollector initializes the collectors and assigns fqName and help description for exported metrics
func newUsageCollector() *UsageCollector {
	return &UsageCollector{
		callerIDLookups:         prometheus.NewDesc("twil_callerIDLookups", "Total CallerID Lookups", nil, nil),
		calls:                   prometheus.NewDesc("twil_calls", "Total Call Minutes", nil, nil),
		callsClient:             prometheus.NewDesc("twil_calls_client", "Total Client Call Minutes", nil, nil),
		callsSip:                prometheus.NewDesc("twil_calls_sip", "SIP Minutes", nil, nil),
		callsInbound:            prometheus.NewDesc("twil_calls_inbound", "Inbound Voice Minutes", nil, nil),
		callsInboundLocal:       prometheus.NewDesc("twil_calls_inbound_local", "Inbound Local Calls", nil, nil),
		callsInboundMobile:      prometheus.NewDesc("twil_calls_mobile", "Inbound Mobile Calls", nil, nil),
		callsInboundTollFree:    prometheus.NewDesc("twil_calls_tollfree", "Inbound Toll Free Calls", nil, nil),
		callsOutbound:           prometheus.NewDesc("twil_calls_outbound", "Outbound Voice Minutes", nil, nil),
		phoneNumbers:            prometheus.NewDesc("twil_phonenumbers", "Phone Numbers", nil, nil),
		phoneNumbersMobile:      prometheus.NewDesc("twil_phonenumbers_mobile", "Mobile Phone Numbers", nil, nil),
		phoneNumbersLocal:       prometheus.NewDesc("twil_phonenumbers_local", "Local Phone Numbers", nil, nil),
		phoneNumbersTollFree:    prometheus.NewDesc("twil_phonenumbers_tollfree", "Toll Free Phone Numbers", nil, nil),
		shortCodes:              prometheus.NewDesc("twil_shortcodes", "Short Codes", nil, nil),
		shortCodesCustomerOwned: prometheus.NewDesc("twil_shortcodes_customer_owned", "Customer Owned Short Codes", nil, nil),
		shortCodesRandom:        prometheus.NewDesc("twil_shortcodes_random", "Random Short Codes", nil, nil),
		shortCodesVanity:        prometheus.NewDesc("twil_shortcodes_vanity", "Vanity Short Codes", nil, nil),
		sms:                     prometheus.NewDesc("twil_sms", "SMS", nil, nil),
		smsInbound:              prometheus.NewDesc("twil_sms_inbound", "Inbound SMS", nil, nil),
		smsInboundLongCode:      prometheus.NewDesc("twil_sms_inbound_standard", "Standard Inbound SMS", nil, nil),
		smsInboundShortCode:     prometheus.NewDesc("twil_sms_inbound_shortcode", "Short Code Inbound SMS", nil, nil),
		smsOutbound:             prometheus.NewDesc("twil_sms_outbound", "Outbound SMS", nil, nil),
		smsOutboundLongCode:     prometheus.NewDesc("twil_sms_outbound_standard", "Standard Outbound SMS", nil, nil),
		smsOutboundShortCode:    prometheus.NewDesc("twil_sms_outbound_shortcode", "Short Code Outbound SMS", nil, nil),
		mms:                     prometheus.NewDesc("twil_mms", "MMS", nil, nil),
		mmsInbound:              prometheus.NewDesc("twil_mms_inbound", "Inbound MMS", nil, nil),
		mmsInboundLongCode:      prometheus.NewDesc("twil_mms_inbound_standard", "Standard Inbound MMS", nil, nil),
		mmsInboundShortCode:     prometheus.NewDesc("twil_mms_inbound_shortcode", "Short Code Inbound MMS", nil, nil),
		mmsOutbound:             prometheus.NewDesc("twil_mms_outbound", "Outbound MMS", nil, nil),
		mmsOutboundLongCode:     prometheus.NewDesc("twil_mms_outbound_standard", "Standard Outbound MMS", nil, nil),
		mmsOutboundShortCode:    prometheus.NewDesc("twil_mms_outbound_shortcode", "Short Code Outbound MMS", nil, nil),
		recordings:              prometheus.NewDesc("twil_recordings", "Recordings", nil, nil),
		recordingsStorage:       prometheus.NewDesc("twil_recordings_storage", "Recordings Storage", nil, nil),
		transcriptions:          prometheus.NewDesc("twil_transcriptions", "Transcriptions", nil, nil),
		mediaStorage:            prometheus.NewDesc("twil_mediastorage", "Media Storage", nil, nil),
		authySMSOutbound:        prometheus.NewDesc("twil_authy_sms_outbound", "Authy/Verify Outbound SMS Messages", nil, nil),
		authyCallsOutbound:      prometheus.NewDesc("twil_authy_calls_outbound", "Authy/Verify Outbound Calls", nil, nil),
		authyAuthentications:    prometheus.NewDesc("twil_authy_authentications", "Authy Authentications", nil, nil),
		authyPhoneVerifications: prometheus.NewDesc("twil_authy_phone_verifications", "Verify", nil, nil),
		authyPhoneIntelligence:  prometheus.NewDesc("twil_authy_phone_intelligence", "Authy Phone Intelligence Requests", nil, nil),
		authyMonthlyFees:        prometheus.NewDesc("twil_authy_monthly_fees", "Authy Monthly Fees", nil, nil),
		monitorStorage:          prometheus.NewDesc("twil_monitor_storage", "Monitor Events Storage", nil, nil),
		monitorReads:            prometheus.NewDesc("twil_monitor_reads", "Monitor Events API Reads", nil, nil),
		monitorWrites:           prometheus.NewDesc("twil_monitor_writes", "Monitor Events API Writes", nil, nil),
		taskRouterTasks:         prometheus.NewDesc("twil_task_router_tasks", "Task Router Tasks Created", nil, nil),
		turnMegabytes:           prometheus.NewDesc("twil_turn_megabytes", "TURN Megabytes", nil, nil),
		callRecordings:          prometheus.NewDesc("twil_call_recordings", "Call Recordings", nil, nil),
		trunkingRecordings:      prometheus.NewDesc("twil_trunking_recordings", "Trunking Recordings", nil, nil),
		trunkingTermination:     prometheus.NewDesc("twil_trunking_termination", "Trunking Termination Minutes", nil, nil),
		trunkingOrigination:     prometheus.NewDesc("twil_trunking_origination", "Trunking Origination Minutes", nil, nil),
	}
}

//Describe initializes channels used to pull Metrics
func (c *UsageCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.callerIDLookups
	ch <- c.calls
	ch <- c.callsClient
	ch <- c.callsSip
	ch <- c.callsInbound
	ch <- c.callsInboundLocal
	ch <- c.callsInboundMobile
	ch <- c.callsInboundTollFree
	ch <- c.callsOutbound
	ch <- c.phoneNumbers
	ch <- c.phoneNumbersMobile
	ch <- c.phoneNumbersLocal
	ch <- c.phoneNumbersTollFree
	ch <- c.shortCodes
	ch <- c.shortCodesCustomerOwned
	ch <- c.shortCodesRandom
	ch <- c.shortCodesVanity
	ch <- c.sms
	ch <- c.smsInbound
	ch <- c.smsInboundLongCode
	ch <- c.smsInboundShortCode
	ch <- c.smsOutbound
	ch <- c.smsOutboundLongCode
	ch <- c.smsOutboundShortCode
	ch <- c.mms
	ch <- c.mmsInbound
	ch <- c.mmsInboundLongCode
	ch <- c.mmsInboundShortCode
	ch <- c.mmsOutbound
	ch <- c.mmsOutboundLongCode
	ch <- c.mmsOutboundShortCode
	ch <- c.recordings
	ch <- c.recordingsStorage
	ch <- c.transcriptions
	ch <- c.mediaStorage
	ch <- c.authySMSOutbound
	ch <- c.authyCallsOutbound
	ch <- c.authyAuthentications
	ch <- c.authyPhoneVerifications
	ch <- c.authyPhoneIntelligence
	ch <- c.authyMonthlyFees
	ch <- c.monitorStorage
	ch <- c.monitorReads
	ch <- c.monitorWrites
	ch <- c.taskRouterTasks
	ch <- c.turnMegabytes
	ch <- c.callRecordings
	ch <- c.trunkingRecordings
	ch <- c.trunkingTermination
	ch <- c.trunkingOrigination
}

//Collect gathers the metrics
func (c *UsageCollector) Collect(ch chan<- prometheus.Metric) {

	client := http.Client{}

	reqURL := "https://api.twilio.com/2010-04-01/Accounts/" + *Account + "/Usage/Records.json"
	method := "GET"

	req, err := http.NewRequest(method, reqURL, nil)

	if err != nil {
		fmt.Println(err)
	}

	formattedToken := "Basic " + *Token
	req.Header.Add("Authorization", formattedToken)
	req.Header.Add("User-Agent", "twil")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var bodyObject Usage
	json.Unmarshal(body, &bodyObject)

	for k := range bodyObject.UsageRecords {
		switch {
		case bodyObject.UsageRecords[k].Category == "callerIDLookups":
			ch <- prometheus.MustNewConstMetric(c.callerIDLookups, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "calls":
			ch <- prometheus.MustNewConstMetric(c.calls, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "calls-client":
			ch <- prometheus.MustNewConstMetric(c.callsClient, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "calls-sip":
			ch <- prometheus.MustNewConstMetric(c.callsSip, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "calls-inbound":
			ch <- prometheus.MustNewConstMetric(c.callsInbound, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "calls-inbound-local":
			ch <- prometheus.MustNewConstMetric(c.callsInboundLocal, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "calls-inbound-mobile":
			ch <- prometheus.MustNewConstMetric(c.callsInboundMobile, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "calls-inbound-tollfree":
			ch <- prometheus.MustNewConstMetric(c.callsInboundTollFree, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "calls-outbound":
			ch <- prometheus.MustNewConstMetric(c.callsOutbound, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "phonenumbers":
			ch <- prometheus.MustNewConstMetric(c.phoneNumbers, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "phonenumbers-mobile":
			ch <- prometheus.MustNewConstMetric(c.phoneNumbersMobile, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "phonenumbers-local":
			ch <- prometheus.MustNewConstMetric(c.phoneNumbersLocal, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "phonenumbers-tollfree":
			ch <- prometheus.MustNewConstMetric(c.phoneNumbersTollFree, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "shortcodes":
			ch <- prometheus.MustNewConstMetric(c.shortCodes, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "shortcodes-customerowned":
			ch <- prometheus.MustNewConstMetric(c.shortCodesCustomerOwned, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "shortcodes-random":
			ch <- prometheus.MustNewConstMetric(c.shortCodesRandom, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "shortcodes-vanity":
			ch <- prometheus.MustNewConstMetric(c.shortCodesVanity, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "sms":
			ch <- prometheus.MustNewConstMetric(c.sms, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "sms-inbound":
			ch <- prometheus.MustNewConstMetric(c.smsInbound, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "sms-inbound-longcode":
			ch <- prometheus.MustNewConstMetric(c.smsInboundLongCode, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "sms-inbound-shortcode":
			ch <- prometheus.MustNewConstMetric(c.smsInboundShortCode, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "sms-outbound":
			ch <- prometheus.MustNewConstMetric(c.smsOutbound, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "sms-outbound-longcode":
			ch <- prometheus.MustNewConstMetric(c.smsOutboundLongCode, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "sms-outbound-shortcode":
			ch <- prometheus.MustNewConstMetric(c.smsOutboundShortCode, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "mms":
			ch <- prometheus.MustNewConstMetric(c.mms, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "mms-inbound":
			ch <- prometheus.MustNewConstMetric(c.mmsInbound, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "mms-inbound-longcode":
			ch <- prometheus.MustNewConstMetric(c.mmsInboundLongCode, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "mms-inbound-shortcode":
			ch <- prometheus.MustNewConstMetric(c.mmsInboundShortCode, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "mms-outbound":
			ch <- prometheus.MustNewConstMetric(c.mmsOutbound, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "mms-outbound-longcode":
			ch <- prometheus.MustNewConstMetric(c.mmsOutboundLongCode, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "mms-outbound-shortcode":
			ch <- prometheus.MustNewConstMetric(c.mmsOutboundShortCode, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "recordings":
			ch <- prometheus.MustNewConstMetric(c.recordings, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "recordingstorage":
			ch <- prometheus.MustNewConstMetric(c.recordingsStorage, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "transcriptions":
			ch <- prometheus.MustNewConstMetric(c.transcriptions, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "mediastorage":
			ch <- prometheus.MustNewConstMetric(c.mediaStorage, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "authy-sms-outbound":
			ch <- prometheus.MustNewConstMetric(c.authySMSOutbound, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "authy-calls-outbound":
			ch <- prometheus.MustNewConstMetric(c.authyCallsOutbound, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "authy-authentications":
			ch <- prometheus.MustNewConstMetric(c.authyAuthentications, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "authy-phone-verifications":
			ch <- prometheus.MustNewConstMetric(c.authyPhoneVerifications, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "authy-phone-intelligence":
			ch <- prometheus.MustNewConstMetric(c.authyPhoneIntelligence, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "authy-monthly-fees":
			ch <- prometheus.MustNewConstMetric(c.authyMonthlyFees, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "monitor-storage":
			ch <- prometheus.MustNewConstMetric(c.monitorStorage, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "monitor-reads":
			ch <- prometheus.MustNewConstMetric(c.monitorReads, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "monitor-write":
			ch <- prometheus.MustNewConstMetric(c.monitorWrites, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "taskrouter-tasks":
			ch <- prometheus.MustNewConstMetric(c.taskRouterTasks, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "turnmegabytes":
			ch <- prometheus.MustNewConstMetric(c.turnMegabytes, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "calls-recordings":
			ch <- prometheus.MustNewConstMetric(c.callRecordings, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "trunking-recordings":
			ch <- prometheus.MustNewConstMetric(c.trunkingRecordings, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "trunking-termination":
			ch <- prometheus.MustNewConstMetric(c.trunkingTermination, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		case bodyObject.UsageRecords[k].Category == "trunking-origination":
			ch <- prometheus.MustNewConstMetric(c.trunkingOrigination, prometheus.CounterValue, bodyObject.UsageRecords[k].Count)
		}
	}

}

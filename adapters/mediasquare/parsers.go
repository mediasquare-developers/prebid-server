package mediasquare

import (
	"encoding/json"
	"fmt"

	"github.com/prebid/openrtb/v20/openrtb2"
)

// parserDSA: Struct used to extracts dsa content of a json.
type parserDSA struct {
	DSA interface{} `json:"dsa,omitempty"`
}

// setContent: Unmarshal a []byte into the parserDSA struct.
func (parser *parserDSA) setContent(extJsonBytes []byte) error {
	if len(extJsonBytes) > 0 {
		if err := json.Unmarshal(extJsonBytes, parser); err != nil {
			return errorWritter("<setContent(*parserDSA)> extJsonBytes", err, false)
		}
		return nil
	}
	return errorWritter("<setContent(*parserDSA)> extJsonBytes", nil, true)
}

// getValue: Returns the DSA value as a string, defaultly returns empty-string.
func (parser parserDSA) getValue(request *openrtb2.BidRequest) string {
	if request == nil || request.Regs == nil {
		return ""
	}
	parser.setContent(request.Regs.Ext)
	if parser.DSA != nil {
		return fmt.Sprint(parser.DSA)
	}
	return ""
}

// parserGDPR: Struct used to extract pair of GDPR/Consent of a json.
type parserGDPR struct {
	GDPR    interface{} `json:"gdpr,omitempty"`
	Consent interface{} `json:"consent,omitempty"`
}

// setContent: Unmarshal a []byte into the parserGDPR struct.
func (parser *parserGDPR) setContent(extJsonBytes []byte) error {
	if len(extJsonBytes) > 0 {
		if err := json.Unmarshal(extJsonBytes, parser); err != nil {
			return errorWritter("<setContent(*parserGDPR)> extJsonBytes", err, false)
		}
		return nil
	}
	return errorWritter("<setContent(*parserGDPR)> extJsonBytes", nil, true)
}

// value: Returns the consent or GDPR-string depending of the parserGDPR content, defaulty return empty-string.
func (parser *parserGDPR) value() string {
	switch {
	case parser.Consent != nil:
		return fmt.Sprint(parser.Consent)
	case parser.GDPR != nil:
		return fmt.Sprint(parser.GDPR)
	}
	return ""
}

// getValue: Returns the consent or GDPR-string depending on the openrtb2.User content, defaultly returns empty-string.
func (parser parserGDPR) getValue(field string, request *openrtb2.BidRequest) string {
	if request == nil {
		return ""
	}
	switch {
	case field == "consent_requirement" && request.Regs != nil:
		if ptrInt8ToBool(request.Regs.GDPR) {
			return "true"
		}
		return "false"
	case field == "consent_string" && request.User != nil:
		if len(request.User.Consent) > 0 {
			return request.User.Consent
		}
		parser.setContent(request.User.Ext)
		return parser.value()
	default:
		return ""
	}
}
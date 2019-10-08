package parser

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func Test_cyprusPost_Parse(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		track      string
		wantEvents []Event
		wantErr    bool
	}{
		{
			name: "simple",
			response: `
<html>
<head>
</head>
<body>
<table class="table-striped table-bordered" width="100%" border="0" cellspacing="0" cellpadding="3">
  <tbody>
    <tr align="left" class="tabl2">
      <td colspan="6"><strong>&nbsp;ITEM&nbsp;n°&nbsp;CO820479485DE</strong></td>
      <td align="center">&nbsp;</td>
    </tr>
    <tr align="center">
      <td class="tabmen" width="13%">Date &amp; Time/ Ημερομηνία &amp; Ώρα </td>
      <td class="tabmen" width="13%">Country / Χώρα </td>
      <td class="tabmen" width="13%">Location / Τοποθεσία </td>
      <td class="tabmen" width="22%">Description / Περιγραφή</td>
      <td class="tabmen" width="15%">Next Office / Επόμενος σταθμός</td>
      <td class="tabmen" width="27%">Extra Information / Επιπρόσθετες πληροφορίες</td>
    </tr>
    

    <tr class="tabl1">
      <td align="center" width="13%">3/4/2019 3:51:00 PM</td>
      <td align="center" width="13%">Germany</td>
      <td align="center" width="13%">DE-63179</td>
      <td align="center" width="22%">Receive item from sender/ Κατάθεση αντικειμένου από τον αποστολέα</td>
      <td align="center" width="15%"> </td>
      <td align="center" width="27%"></td>
  
    </tr>
  </tbody>
</table>
</body>
</html>
`,
			track: "",
			wantEvents: []Event{
				{
					When: time.Date(2019, 3, 4, 15, 51, 00, 0, time.UTC),
					Description: []string{
						"Germany",
						"DE-63179",
						"Receive item from sender/ Κατάθεση αντικειμένου από τον αποστολέα",
						" ",
						"",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(tt.response))
			}))
			c := NewCyprusPost()
			c.PageUrl = s.URL + "?track=%s"
			gotEvents, err := c.Parse(tt.track)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotEvents, tt.wantEvents) {
				t.Errorf("Parse() gotEvents = %v, want %v", gotEvents, tt.wantEvents)
			}
		})
	}
}

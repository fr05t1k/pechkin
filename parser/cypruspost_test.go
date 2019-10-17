package parser

import (
	"github.com/fr05t1k/pechkin/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_cyprusPost_Parse(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		track      string
		wantEvents []storage.Event
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
			wantEvents: []storage.Event{
				{
					EventAt:     time.Date(2019, 3, 4, 15, 51, 00, 0, time.UTC),
					Description: "Germany\nDE-63179\nReceive item from sender/ Κατάθεση αντικειμένου από τον αποστολέα\n \n\n",
				},
			},
			wantErr: false,
		},
		{
			name: "complex",
			response: `
<html>
<body leftmargin="0" topmargin="0" marginwidth="0" marginheight="0">

 

<table class="table-striped table-bordered" width="100%" border="0" cellspacing="0" cellpadding="3">
  <tbody>
    <tr align="left" class="tabl2">
      <td colspan="6"><strong>&nbsp;ITEM&nbsp;n°&nbsp;E</strong></td>
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
      <td align="center" width="13%">10/25/2013 3:12:19 PM</td>
      <td align="center" width="13%">Cyprus</td>
      <td align="center" width="13%">LEFKOSIA D.P.O. COUNTER 1901</td>
      <td align="center" width="22%">Receive item from sender/ Κατάθεση αντικειμένου από τον αποστολέα</td>
      <td align="center" width="15%"> </td>
      <td align="center" width="27%"></td>
  
    </tr>
    

    <tr class="tabl2">
      <td align="center" width="13%">11/22/2013 12:44:01 PM</td>
      <td align="center" width="13%">Cyprus</td>
      <td align="center" width="13%">LARNAKA D.P.O. 6900</td>
      <td align="center" width="22%">Send item to next processing point/ Αποστολή αντικειμένου στο επόμενο σημείο διαχείρισης. </td>
      <td align="center" width="15%">PERIVOLIA LARNAKAS 7560</td>
      <td align="center" width="27%"></td>
  
    </tr>
    

    <tr class="tabl1">
      <td align="center" width="13%">11/29/2013 12:46:59 PM</td>
      <td align="center" width="13%">Cyprus</td>
      <td align="center" width="13%">LEMESOS D.P.O. 3900</td>
      <td align="center" width="22%">Send item to next processing point/ Αποστολή αντικειμένου στο επόμενο σημείο διαχείρισης. </td>
      <td align="center" width="15%">LEMESOS B.O.2 3904 -LEOFOROS MAKARIOU</td>
      <td align="center" width="27%"></td>
  
    </tr>
    

    <tr class="tabl2">
      <td align="center" width="13%">2/19/2014 1:45:44 PM</td>
      <td align="center" width="13%">Cyprus</td>
      <td align="center" width="13%">PAFOS D.P.O. 8900</td>
      <td align="center" width="22%">Send item to next processing point/ Αποστολή αντικειμένου στο επόμενο σημείο διαχείρισης. </td>
      <td align="center" width="15%">PEGEIA 8906</td>
      <td align="center" width="27%"></td>
  
    </tr>
    

    <tr class="tabl1">
      <td align="center" width="13%">8/18/2014 8:17:00 AM</td>
      <td align="center" width="13%">Cyprus</td>
      <td align="center" width="13%">LAKATAMEIA 1915</td>
      <td align="center" width="22%">Receive item at delivery office (Inb)/ Παραλαβή αντικειμένου στο γραφείο παράδοσης</td>
      <td align="center" width="15%">-</td>
      <td align="center" width="27%"></td>
  
    </tr>
    

    <tr class="tabl2">
      <td align="center" width="13%">3/2/2016 7:44:55 AM</td>
      <td align="center" width="13%">Cyprus</td>
      <td align="center" width="13%">LEFKOSIA S.B.O. 1903 - PLATEIA ELEFTERHIAS</td>
      <td align="center" width="22%">Send item to processing point (Otb)/ Αποστολή αντικειμένου προς σημείο διαχείρισης</td>
      <td align="center" width="15%">LARNAKA O/E AIR 6911</td>
      <td align="center" width="27%"></td>
  
    </tr>
    

    <tr class="tabl1">
      <td align="center" width="13%">1/19/2018 1:19:00 PM</td>
      <td align="center" width="13%">Cyprus</td>
      <td align="center" width="13%">LEFKOSIA B.O.2 1906 - PALLOURIOTISSA</td>
      <td align="center" width="22%">Receive item at delivery office (Inb)/ Παραλαβή αντικειμένου στο γραφείο παράδοσης</td>
      <td align="center" width="15%">-</td>
      <td align="center" width="27%"></td>
  
    </tr>
    

    <tr class="tabl2">
      <td align="center" width="13%">4/20/2018 9:59:17 AM</td>
      <td align="center" width="13%">Cyprus</td>
      <td align="center" width="13%">LEMESOS B.O.10 3912 - AMATHOUNTA</td>
      <td align="center" width="22%">Receive item at delivery office (Inb)/ Παραλαβή αντικειμένου στο γραφείο παράδοσης</td>
      <td align="center" width="15%">-</td>
      <td align="center" width="27%"></td>
  
    </tr>
    

    <tr class="tabl1">
      <td align="center" width="13%">11/9/2018 1:11:51 PM</td>
      <td align="center" width="13%">Cyprus</td>
      <td align="center" width="13%">LEMESOS B.O.10 3912 - AMATHOUNTA</td>
      <td align="center" width="22%">Deliver item to Addressee (Inb)/ Παράδοση αντικειμένου στον παραλήπτη
</td>
      <td align="center" width="15%">-</td>
      <td align="center" width="27%">SIGNED BY : GUY PHILIPS</td>
  
    </tr>
    

    <tr class="tabl2">
      <td align="center" width="13%">11/27/2018 4:56:00 PM</td>
      <td align="center" width="13%">Denmark</td>
      <td align="center" width="13%">FREDERICIA PARCEL</td>
      <td align="center" width="22%">Receive item at country of destination/ Παραλαβή αντικειμένου στη χώρα προορισμού</td>
      <td align="center" width="15%">-</td>
      <td align="center" width="27%"></td>
  
    </tr>
    
  </tbody>
</table>

<br>
<br>


</body></html>
`,
			track: "",
			wantEvents: []storage.Event{
				{
					EventAt:     time.Date(2013, 10, 25, 15, 12, 19, 0, time.UTC),
					Description: "Cyprus\nLEFKOSIA D.P.O. COUNTER 1901\nReceive item from sender/ Κατάθεση αντικειμένου από τον αποστολέα\n \n\n",
				},
				{
					EventAt:     time.Date(2013, 11, 22, 12, 44, 01, 0, time.UTC),
					Description: "Cyprus\nLARNAKA D.P.O. 6900\nSend item to next processing point/ Αποστολή αντικειμένου στο επόμενο σημείο διαχείρισης. \nPERIVOLIA LARNAKAS 7560\n\n",
				},
				{
					EventAt:     time.Date(2013, 11, 29, 12, 46, 59, 0, time.UTC),
					Description: "Cyprus\nLEMESOS D.P.O. 3900\nSend item to next processing point/ Αποστολή αντικειμένου στο επόμενο σημείο διαχείρισης. \nLEMESOS B.O.2 3904 -LEOFOROS MAKARIOU\n\n",
				},
				{
					EventAt:     time.Date(2014, 2, 19, 13, 45, 44, 0, time.UTC),
					Description: "Cyprus\nPAFOS D.P.O. 8900\nSend item to next processing point/ Αποστολή αντικειμένου στο επόμενο σημείο διαχείρισης. \nPEGEIA 8906\n\n",
				},
				{
					EventAt:     time.Date(2014, 8, 18, 8, 17, 00, 0, time.UTC),
					Description: "Cyprus\nLAKATAMEIA 1915\nReceive item at delivery office (Inb)/ Παραλαβή αντικειμένου στο γραφείο παράδοσης\n-\n\n",
				},
				{
					EventAt:     time.Date(2016, 3, 2, 7, 44, 55, 0, time.UTC),
					Description: "Cyprus\nLEFKOSIA S.B.O. 1903 - PLATEIA ELEFTERHIAS\nSend item to processing point (Otb)/ Αποστολή αντικειμένου προς σημείο διαχείρισης\nLARNAKA O/E AIR 6911\n\n",
				},
				{
					EventAt:     time.Date(2018, 1, 19, 13, 19, 0, 0, time.UTC),
					Description: "Cyprus\nLEFKOSIA B.O.2 1906 - PALLOURIOTISSA\nReceive item at delivery office (Inb)/ Παραλαβή αντικειμένου στο γραφείο παράδοσης\n-\n\n",
				},
				{
					EventAt:     time.Date(2018, 4, 20, 9, 59, 17, 0, time.UTC),
					Description: "Cyprus\nLEMESOS B.O.10 3912 - AMATHOUNTA\nReceive item at delivery office (Inb)/ Παραλαβή αντικειμένου στο γραφείο παράδοσης\n-\n\n",
				},
				{
					EventAt:     time.Date(2018, 11, 9, 13, 11, 51, 0, time.UTC),
					Description: "Cyprus\nLEMESOS B.O.10 3912 - AMATHOUNTA\nDeliver item to Addressee (Inb)/ Παράδοση αντικειμένου στον παραλήπτη\n\n-\nSIGNED BY : GUY PHILIPS\n",
				},
				{
					EventAt:     time.Date(2018, 11, 27, 16, 56, 0, 0, time.UTC),
					Description: "Denmark\nFREDERICIA PARCEL\nReceive item at country of destination/ Παραλαβή αντικειμένου στη χώρα προορισμού\n-\n\n",
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
			if !assert.Len(t, gotEvents, len(tt.wantEvents)) {
				return
			}
			if !assert.Equal(t, gotEvents, tt.wantEvents) {
				return
			}
		})
	}
}

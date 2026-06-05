package models

//maqolaga tegishli fayillarni saqlash uchun alohida structura

type Author struct{

        Name            string `json:"name"`//ismi sharifi
        Affiliation     string `json:"affiliation"`//ish joyi
        ORCID           string `json:"orcid"`// halqaro mualiflar id si 
}

type Article struct{
	     ID             string `json:"id"`
		 Title          string `json:"title"` //1 mavzu uzb rus eng
		 AccessType     string `json:"accessType"` //2foydalnuvchi huquqi (ochiq/yopiq)
		 Authors        []Author `json:"authors"` //3 muallaiflar royhati
		 Abstract       string `json:"abstract"`//4anotatsiyasi

		 KeyWords       []string `json:"keyWords"` //5 kalit so'zlar
		 Journal        string `json:"journal"`//6 jurnal nomi
		 Publisher      string `json:"publisher"`//7 jurnal chiqaruvchi tashkilot
		 PublisherDate  string `json:"publisherDate"`//8 chop etilgan sana
		 

		 DOI            string `json:"doi"`//9 maqola DOI
		 URL            string `json:"url"`//10 maqola url manzil
		 PDFUrl         string `json:"pdfUrl"`//11 pdf fayl manzili
		 SourceURL      string `json:"sourceUrl"`//12 manba URL	
}

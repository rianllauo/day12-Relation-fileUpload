// megambil setiap input dari id nya
let username = document.getElementById('username'),
   email = document.getElementById('email'),
   phone = document.getElementById('phone'),
   subject = document.getElementById('subject'),
   message = document.getElementById('message')
   button = document.getElementById('submit'),      
   errorMessage = document.getElementsByClassName('errorMessage')

// membuat function ketika button di klik
function onSubmit(){
   // menjalankan fungsi submit data pada setiap input untuk
   // mengecek apakah semua input tidak kosong  
   // dengan parameter id dari setiap input, index, dan pesan errornya
   submitData(username, 0, 'Please input your name')
   submitData(email, 1, 'Please input your email')
   submitData(phone, 2, 'Please input your phone number')
   submitData(subject, 3, 'Please input your subject')
   submitData(message, 4, 'Please input your message')  
}

// membuat function saat button di klik untuk mensubmit data
function submitData (id, index, warning) {

   // fungsi untuk menampilkan pesan error saat input tidak di isi
   function fail() {
      errorMessage[index].innerHTML = warning;
      id.style.border = "2.5px solid #e11d48"
   }

   // melakukan pengecekan saat input tidak di isi
   if(id.value === ""){
      return fail()
   } 

   // pengecekan kembali saat semua inputan sudah di isi / tidak kosong
   if(username.value && email.value && phone.value && subject.value && message.value !== ""){
     return emailTo()
   }
   
   // menjalankan fungsi untuk mengirim email
   function emailTo(){
      let emailReciver = "rianllauo@gmail.com"
      let link = document.createElement('a')
      link.href = `mailto: ${emailReciver}?subject=${subject.value}&body=Hallo nama saya ${username.value}, ${message.value}, silahkan kontak saya di nomer ${phone.value}`
      link.click()
   }
}

// membuat fungsi untuk menghilangkan error saat 
// input di masukan
function onInput(){
   inputFill(username, 0);
   inputFill(email, 1);
   inputFill(phone, 2);
   inputFill(subject, 3);
   inputFill(message, 4);
}
function inputFill(id, index){
   // membuat function success
   // untuk menghilangkan error message
   function success() {
      errorMessage[index].innerHTML = "";
      id.style.border = "2.5px solid #6366f1";
   }

   // melakukan cek terhadap value dari inputan 
   // apa sudah di isi atau belum
   if(id.value !== ""){
      return success()
   }
}
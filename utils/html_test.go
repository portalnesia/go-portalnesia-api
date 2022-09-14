package utils

import (
	"regexp"
	"testing"
)

const newsText string = `<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN" "http://www.w3.org/TR/REC-html40/loose.dtd">
<center><a data-src="https://cdn-asset.jawapos.com/wp-content/uploads/2019/04/Napi-640x351.jpg" data-caption="Pelaksanaan Permenkumham ini merupakan langkah yang ditempuh untuk melindungi hak kesehatan WBP di masa pandemi Covid-19" data-fancybox="images" data-options='{"protect" : "true" }' style="cursor:pointer"><picture><source type="image/webp" data-srcset="https://content.portalnesia.com/img?size=400&amp;output=webp&amp;type=url&amp;image=https%3A%2F%2Fcdn-asset.jawapos.com%2Fwp-content%2Fuploads%2F2019%2F04%2FNapi-640x351.jpg"></source><source type="image/jpeg" data-srcset="https://content.portalnesia.com/img?size=400&amp;type=url&amp;image=https%3A%2F%2Fcdn-asset.jawapos.com%2Fwp-content%2Fuploads%2F2019%2F04%2FNapi-640x351.jpg"></source><img data-src="https://content.portalnesia.com/img?size=400&amp;type=url&amp;image=https%3A%2F%2Fcdn-asset.jawapos.com%2Fwp-content%2Fuploads%2F2019%2F04%2FNapi-640x351.jpg" loading="lazy" style="width:85%;" class="text-center my-3 lazyload"></picture></a></center><html><body><p>JawaPos.com &ndash; Kementerian Hukum dan Hak Asasi Manusia (Kemenkumham) kembali memperpanjang program pemberian hak Integrasi dan Asimilasi di rumah bagi narapidana dan Anak sebagai pencegahan dan penanggulangan penyebaran Coronavirus Disease (Covid-19). Hal tersebut diwujudkan dengan dikeluarkannya Peraturan Menteri Hukum dan Hak Asasi Manusia Republik Indonesia (Permenkumham RI) Nomor 43 Tahun 2021.</p>
<p>Adapun Permenkumham ini merupakan Perubahan Kedua atas Permenkumham RI Nomor 32 Tahun 2020 dan Permenkumham RI Nomor 24 Tahun 2021 tentang Syarat dan Tata Cara Pemberian Asimilasi, PB, CMB, dan CB bagi Narapidana dan Anak dalam rangka Pencegahan dan Penanggulangan Penyebaran Covid-19.</p>
<p>Kepala Bagian Humas dan Protokol Direktorat Jenderal Pemasyarakatan, Rika Aprianti mengungkapkan, hal ini merupakan upaya lanjutan Kemenkumham dalam mencegah dan menanggulangi penyebaran Covid-19 di Lembaga Pemasyarakatan, Rumah Tahanan Negara, dan Lembaga Pembinaan Khusus Anak melalui pemberian Asimilasi dan Integrasi.</p>
<p>&ldquo;Pelaksanaan Permenkumham ini merupakan langkah yang ditempuh untuk melindungi hak kesehatan WBP di masa pandemi Covid-19 yang telah terjadi sejak awal tahun 2020, terlebih saat ini muncul berbagai varian baru yang harus kita waspadai,&rdquo; kata Rika dalam keterangannya, Minggu (2/1).[IKLAN_IKLAN]</p>
<p>Penerbitan Permenkumham tersebut menjadi respon terhadap pandemi yang masih berlangsung di berbagai belahan dunia hingga saat ini. Untuk itu, Rika kembali menegaskan bahwa Pemasyarakatan akan melaksanakan ketentuan tata cara pemberian Asimilasi, Pembebasan Bersyarat, Cuti Menjelang Bebas, dan Cuti Bersyarat sesuai peraturan yang ada.</p>
<p>&ldquo;Adapun perubahan yang dilakukan terkait narapidana penerima Asimilasi dan perluasan jangkauan penerima hak Integrasi dan Asimilasi bagi narapidana dan Anak. Bila semula berlaku bagi narapidana yang 2/3 masa pidana dan Anak yang 1/2 masa pidananya hingga 31 Desember 2021, kini diperpanjang hingga 30 Juni 2022,&rdquo; papar Rika.</p>
<p>Terkait pelaksanaan Permenkumham RI tersebut, pihaknya menekankan seluruh proses layanan Asimilasi dan Integrasi tidak dipungut biaya apa pun. Karena itu, seluruh petugas perlu mencermati dan melaksanakan peraturan ini agar tidak terjadi kesalahan.</p>
<p>&ldquo;Nantinya akan makin banyak yang melaksanakan Asimilasi dan Integrasinya di rumah, tentunya dengan pengawasan dari Pembimbing Kemasyarakatan di Balai Pemasyarakatan. Kami juga berharap masyarakat mau berperan serta mengawasi dan mendukung pelaksanaan Asimilasi di rumah. Kami akan terus melakukan upaya pencegahan, penanggulangan, dan penanganan penanganan penyebaran Covid-19 dengan lebih optimal,&rdquo; tegas Rika.[IKLAN_IKLAN_IKLAN]</p>
<p>Sebelumnya, sejak awal pandemi Kemenkumham telah mengeluarkan kebijakan pelaksanaan pemberian Hak Asimilasi dan Integrasi di Rumah melalui Permenkumham RI Nomor 10 Tahun 2020 tentang Pemberian Asimilasi dan Hak Integrasi Bagi Narapidana dan Anak Dalam Rangka Pencegahan dan Penanggulangan Penyebaran Covid-19.</p>
<p>&ldquo;Hingga saat ini kebijakan tersebut telah berhasil merumahkan 94.047 narapidana dan 2.026 Anak untuk menjalankan hak Integrasi dan 115.798 narapidana dan Anak untuk menjalankan hak Asimilasi di rumah,&rdquo; pungkas Rika.</p></body></html>
`

const blogText string = `<p style="text-align: justify;">Sebelumnya kita sudah membahas tentang verifikasi dua langkah dengan menggunakan security key browser (sidik jari, pin, dsb). Kali ini kami mau membahas topik yang sama, tetapi menggunakan metode yang berbeda, yaitu <strong><em>Mobile App Authentication</em></strong>.</p><p style="text-align:justify">Mobile App Authentication adalah metode verifikasi menggunakan kode sekali pakai (one time password). Jadi ketika kamu berhasil login, kamu akan dimintai kode yang kamu dapat dari aplikasi sejenis google authenticator. Kode tersebut unik dan rahasia, <strong>kamu pun tidak boleh memberikan kode tersebut kepada siapapun agar akun kamu aman</strong>.</p><p style="text-align:justify">&nbsp;</p><p style="text-align:justify">Jadi bagaimana cara mengaktifkannya?</p><p style="text-align:justify">&nbsp;</p><h2 style="text-align:justify">Download aplikasi</h2><p style="text-align:justify">Pertama-tama kamu harus mendownload aplikasinya. Kamu bisa menggunakan <a href="https://support.google.com/accounts/answer/1066447?hl=en" target="_blank">Google Authenticator</a>, <a href="https://duo.com/" target="_blank">Duo Mobile</a>, atau aplikasi sejenis lainnya. Pada tutorial ini, kami menggunakan aplikasi duo mobile. Tetapi secara umum, langkah-langkahnya gak jauh beda kok.</p><p style="text-align:justify">&nbsp;</p><h2 style="text-align:justify">Aktifkan&nbsp;Mobile App Authentication</h2><div style="text-align:center"><figure class="image" style="display:inline-block"><img alt="" height="202" src="https://i.imgur.com/EWLZN21.png" width="400" /><figcaption>Halaman setting Portalnesia</figcaption></figure></div><p style="text-align:justify">Aktifkan <strong>Mobile App Authentication</strong> pada <a href="https://portalnesia.com/setting/security">halaman setting</a> Portalnesia. Juga jangan lupa untuk mengaktifkan <strong>Two Factor Authentication</strong>&nbsp;terlebih dahulu.</p><p style="text-align:justify">Masukkan kata sandi kamu, dan klik tombol <strong>Submit</strong>.</p><p style="text-align:justify">&nbsp;</p><h2 style="text-align:justify">Tambahkan Portalnesia pada aplikasi</h2><div style="text-align:center"><figure class="image" style="display:inline-block"><img alt="" height="637" src="https://i.imgur.com/7XMtC5j.jpg" width="319" /><figcaption>Aplikasi Duo Mobile</figcaption></figure></div><p>&nbsp;</p><div style="text-align:center"><figure class="image" style="display:inline-block"><img alt="" height="201" src="https://i.imgur.com/s27XEtP.png" width="400" /><figcaption>Halaman setting Portalnesia</figcaption></figure></div><p>&nbsp;</p><p style="text-align:justify">Tambahkan konfigurasi yang kami berikan, bisa menggunakan QR code, bisa juga dengan menyalin kode dan tempelkan pada aplikasi authenticatormu.</p><p style="text-align:justify">&nbsp;</p><h2 style="text-align:justify">Masukkan kode yang ditampilkan pada aplikasi authenticator</h2><div style="text-align:center"><figure class="image" style="display:inline-block"><img alt="" height="623" src="https://i.imgur.com/2Kx6IoI.jpg" width="312" /><figcaption>Aplikasi Duo Mobile</figcaption></figure></div><p>&nbsp;</p><div style="text-align:center"><figure class="image" style="display:inline-block"><img alt="" height="201" src="https://i.imgur.com/RS1H72H.png" width="400" /><figcaption>Halaman setting Portalnesia</figcaption></figure></div><p>&nbsp;</p><p style="text-align:justify">Sebagai langkah verifikasi bahwa kamu sudah menambahkan konfigurasi yang kami berikan pada aplikasi authenticatormu, kamu harus memasukkan kode yang ditampilkan&nbsp;pada aplikasi authenticator.&nbsp;</p><p style="text-align:justify">Lalu klik tombol <strong>Verify</strong>.</p><p style="text-align:justify">&nbsp;</p><p style="text-align:justify">Yak, kamu telah berhasil menambahkan satu lagi metode verifikasi dua langkah.&nbsp;</p><p style="text-align:justify">&nbsp;</p><p style="text-align:justify"><strong>Lalu bagaimana cara kerja pada saat kamu melakukan&nbsp; proses login?</strong></p><p style="text-align:justify">Caranya mirip dengan langkah-langkah login menggunakan security key</p><div style="text-align:center"><figure class="image" style="display:inline-block"><img alt="" height="428" src="https://i.imgur.com/SJHtiZB.png" width="400" /><figcaption>Verifikasi dua langkah Portalnesia</figcaption></figure></div><p>&nbsp;</p><p style="text-align:justify">Hanya saja pada halaman verifikasi dua langkah, kamu tidak menggunakan sidik jari, tetapi menggunakan 6 kode yang ditampilkan pada aplikasi authenticatormu tadi.</p><p style="text-align:justify">Kode tersebut juga bisa dikirimkan melalui <strong>SMS</strong> (jika kamu sudah melakukan verifikasi nomor HP) dan melalui <strong>Telegram</strong> (jika kamu sudah memverifikasi akun telegram kamu).</p><p style="text-align:justify">&nbsp;</p><p style="text-align:justify">Sekian dulu tutorial kali inii, sampai ketemu lagi :)</p>
`

func TestNewsEncode(t *testing.T) {
	test := NewsEncode(newsText)

	find1 := regexp.MustCompile(`\[IKLAN_IKLAN\]`).FindAllString(test, -1)
	find2 := regexp.MustCompile(`\[IKLAN_IKLAN_IKLAN\]`).FindAllString(test, -1)

	if len(find1) > 0 {
		t.Errorf("NewsEncode [IKLAN_IKLAN] Error, Get %s", find1)
	}
	if len(find2) > 0 {
		t.Errorf("NewsEncode [IKLAN_IKLAN_IKLAN] Error, Get %s", find2)
	}
}

func TestBlogEncode(t *testing.T) {
	test := BlogEncode(blogText, false)

	find1 := regexp.MustCompile(`\<div data\-portalnesia\-action`).FindAllString(test, -1)

	if len(find1) == 0 {
		t.Errorf("BlogEncode Iklan Error, Get %s", find1)
	}
}

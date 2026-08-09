-- =====================================================================
-- SEED DATA: assessment_kpsp_questions + recommended_actions
-- Sumber: mapping_recomendation.pdf
-- Catatan: item "Membantu Orang" (bulan 30) menggunakan domain
-- speech_and_language sesuai literal PDF (bukan hasil koreksi manual).
-- =====================================================================

WITH data AS (
    SELECT DISTINCT * FROM (
        VALUES
    -- =======================================
    -- BULAN 3
    -- =======================================
    (3, 'gross_motor_skills',
    'Pada waktu bayi telentang, apakah masing masing lengan dan tungkai bergerak dengan mudah? Jawab TIDAK bila salah satu atau kedua tungkai atau lengan bayi bergerak tak terarah/tak terkendali.',
    'Gerakan Tubuh',
    'Ajak bayi bergerak bebas saat tengkurap maupun telentang setiap hari agar otot lengan dan tungkainya semakin kuat, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (3, 'socialization',
    'Pada waktu bayi telentang apakah ia melihat dan menatap wajah anda?',
    'Kontak Mata',
    'Seringlah mengajak bayi menatap wajah saat berbicara, menyusui, atau bermain agar kemampuan berinteraksi dan memperhatikan orang di sekitarnya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (3, 'speech_and_language',
    'Apakah bayi dapat mengeluarkan suarasuara lain (ngoceh), disamping menangis?',
    'Merangsang Ocehan',
    'Ajak bayi berbicara, bernyanyi, dan menanggapi setiap suara yang dikeluarkannya agar kemampuan komunikasi awal terus berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (3, 'fine_motor_skills',
    'Pada waktu bayi telentang, apakah ia dapat mengikuti gerakan anda dengan menggerakkan kepalanya dari kanan/kiri ke tengah?',
    'Mengikuti Benda',
    'Gerakkan mainan berwarna cerah secara perlahan ke kanan dan kiri agar bayi belajar mengikuti benda dengan pandangan dan gerakan kepala, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (3, 'fine_motor_skills',
    'Pada waktu bayi telentang, apakah ia dapat mengikuti gerakan anda dengan menggerakkan kepalanya dari satu sisi hampir sampai pada sisi yang lain?',
    'Fokus Mata',
    'Perlihatkan wajah atau mainan pada jarak sekitar 20-30 cm agar bayi terbiasa memusatkan pandangan terhadap objek di depannya, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (3, 'socialization',
    'Pada waktu anda mengajak bayi berbicara dan tersenyum, apakah ia tersenyum kembali kepada anda?',
    'Senyum Sosial',
    'Seringlah mengajak bayi tersenyum, bercanda, dan berbicara dengan ekspresi yang hangat agar bayi mulai merespons interaksi sosial secara alami, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (3, 'gross_motor_skills',
    'Pada waktu bayi telungkup di alas yang datar, apakah ia dapat mengangkat kepalanya?',
    'Belajar Tengkurap',
    'Berikan tummy time beberapa kali sehari di alas yang aman dengan durasi singkat agar kekuatan leher, bahu, dan punggung bayi semakin meningkat, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (3, 'gross_motor_skills',
    'Pada waktu bayi telungkup di alas yang datar, apakah ia dapat mengangkat kepalanya sehingga membentuk sudut 45 derajat?',
    'Angkat Kepala',
    'Letakkan mainan atau wajah Anda di depan bayi saat tengkurap agar ia terdorong mengangkat kepala untuk melihat sekitarnya, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (3, 'gross_motor_skills',
    'Pada waktu bayi telungkup di alas yang datar, apakah ia dapat mengangkat kepalanya dengan tegak?',
    'Kekuatan Leher',
    'Biasakan bayi melakukan tummy time secara bertahap setiap hari untuk membantu memperkuat otot leher sebagai dasar perkembangan gerak berikutnya, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (3, 'speech_and_language',
    'Apakah bayi suka tertawa keras walau tidak digelitik atau diraba-raba?',
    'Merangsang Tawa',
    'Ajak bayi bermain sederhana menggunakan suara, senyuman, atau ekspresi lucu agar ia mulai tertawa dan menikmati interaksi dengan orang di sekitarnya, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    -- =======================================
    -- BULAN 6
    -- =======================================
    (6, 'fine_motor_skills',
    'Pada waktu bayi telentang, apakah ia dapat mengikuti gerakan anda dengan menggerakkan kepala sepenuhnya dari satu sisi ke sisi yang lain?',
    'Mengikuti Gerakan',
    'Gerakkan mainan berwarna cerah secara perlahan ke berbagai arah agar bayi terbiasa mengikuti pergerakan benda menggunakan mata dan kepala, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (6, 'gross_motor_skills',
    'Dapatkah bayi mempertahankan posisi kepala dalam keadaan tegak clan stabil? Jawab TIDAK bila kepala bayi cenderung jatuh ke kanan/kiri atau ke dadanya',
    'Kepala Tegak',
    'Sering ajak bayi duduk dengan penyangga atau digendong dalam posisi tegak agar otot leher dan punggung semakin kuat menopang kepala, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (6, 'fine_motor_skills',
    'Sentuhkan pensil di punggung tangan atau ujung jari bayi. (jangan meletakkan di atas telapak tangan bayi). Apakah bayi dapat menggenggam pensil itu selama beberapa detik?',
    'Belajar Menggenggam',
    'Berikan mainan yang ringan dan mudah digenggam agar bayi terbiasa melatih kekuatan tangan serta koordinasi jari melalui aktivitas bermain, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (6, 'gross_motor_skills',
    'Ketika bayi telungkup di alas datar, apakah ia dapat mengangkat dada dengan kedua lengannya sebagai penyangga?',
    'Mengangkat Dada',
    'Lakukan tummy time setiap hari sambil meletakkan mainan di depan bayi agar ia terdorong menopang tubuh dengan kedua lengannya, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (6, 'speech_and_language',
    'Pernahkah bayi mengeluarkan suara gembira bernada tinggi atau memekik tetapi bukan menangis?',
    'Merangsang Suara',
    'Sering ajak bayi bercakap-cakap, menirukan suara yang dikeluarkannya, dan bernyanyi bersama agar kemampuan bicara awal berkembang lebih baik, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (6, 'gross_motor_skills',
    'Pernahkah bayi berbalik paling sedikit dua kali, dari telentang ke telungkup atau sebaliknya?',
    'Belajar Berguling',
    'Letakkan mainan di samping tubuh bayi untuk mendorongnya berguling ke kanan dan kiri melalui permainan yang aman dan menyenangkan, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (6, 'socialization',
    'Pernahkah anda melihat bayi tersenyurn ketika melihat mainan yang lucu, gambar atau binatang peliharaan pada saat ia bermain sendiri?',
    'Bermain Mandiri',
    'Sediakan mainan yang aman dan menarik agar bayi belajar menikmati waktu bermain sendiri sambil tetap berada dalam pengawasan orang tua, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (6, 'fine_motor_skills',
    'Dapatkah bayi mengarahkan matanya pada benda kecil sebesar kacang, kismis atau uang logam? Jawab TIDAK jika ia tidak dapat mengarahkan matanya.',
    'Fokus Benda',
    'Perlihatkan benda-benda kecil berwarna kontras pada jarak yang aman agar bayi terbiasa memusatkan perhatian dan mengikuti objek dengan penglihatannya, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (6, 'fine_motor_skills',
    'Dapatkah bayi meraih mainan yang diletakkan agak jauh namun masih berada dalam jangkauan tangannya?',
    'Meraih Mainan',
    'Letakkan mainan sedikit di luar jangkauan tangan bayi agar ia terdorong meraih dan mengoordinasikan gerakan tangan serta tubuhnya, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (6, 'gross_motor_skills',
    'Pada posisi bayi telentang, pegang kedua tangannya lalu tarik perlahan-lahan ke posisi duduk. Dapatkah bayi mempertahankan lehernya secara kaku?',
    'Leher Kuat',
    'Biasakan bayi berlatih mengangkat kepala saat ditarik perlahan ke posisi duduk sesuai petunjuk yang benar untuk membantu memperkuat otot leher, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    -- =======================================
    -- BULAN 9
    -- =======================================
    (9, 'gross_motor_skills',
    'Pada posisi bayi telentang, pegang kedua tangannya lalu tarik perlahan-lahan ke posisi duduk. Dapatkah bayi mempertahankan lehernya secara kaku seperti gambar di sebelah kiri ? Jawab TIDAK bila kepala bayi jatuh kembali',
    'Leher Stabil',
    'Latih bayi duduk dan tarik perlahan ke posisi duduk sesuai petunjuk yang benar agar otot leher tetap kuat menopang kepala selama bergerak, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (9, 'fine_motor_skills',
    'Pernahkah anda melihat bayi memindahkan mainan atau kue kering dari satu tangan ke tangan yang lain? Benda-benda panjang seperti sendok atau kerincingan bertangkai tidak ikut dinilai.',
    'Pindah Mainan',
    'Berikan mainan yang mudah digenggam pada salah satu tangan bayi lalu dorong ia memindahkannya ke tangan yang lain melalui aktivitas bermain sehari-hari, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (9, 'fine_motor_skills',
    'Tarik perhatian bayi dengan memperlihatkan selendang, sapu tangan atau serbet, kemudian jatuhkan ke lantai. Apakah bayi mencoba mencarinya? Misalnya mencari di bawah meja atau di belakang kursi?',
    'Mencari Benda',
    'Sembunyikan sebagian mainan favorit di bawah kain atau di balik benda yang mudah dijangkau agar bayi terdorong mencari benda yang hilang dari pandangannya, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (9, 'fine_motor_skills',
    'Apakah bayi dapat memungut dua benda seperti mainan/kue kering, dan masingmasing tangan memegang satu benda pada saat yang sama? Jawab TIDAK bila bayi tidak pernah melakukan perbuatan ini.',
    'Dua Tangan',
    'Ajak bayi bermain menggunakan dua mainan ringan sekaligus agar ia terbiasa memegang satu benda di masing-masing tangan secara bersamaan, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (9, 'gross_motor_skills',
    'Jika anda mengangkat bayi melalui ketiaknya ke posisi berdiri, dapatkah ia menyangga sebagian berat badan dengan kedua kakinya? Jawab YA bila ia mencoba berdiri dan sebagian berat badan tertumpu pada kedua kakinya.',
    'Belajar Berdiri',
    'Biarkan bayi berdiri dengan bantuan pegangan yang kokoh agar otot kaki dan keseimbangannya berkembang secara bertahap, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (9, 'fine_motor_skills',
    'Dapatkah bayi memungut dengan tangannya benda-benda kecil seperti kismis, kacang-kacangan, potongan biskuit, dengan gerakan miring atau menggerapai?',
    'Menjepit Benda',
    'Sediakan makanan ringan berukuran kecil atau mainan yang aman untuk melatih bayi mengambil benda menggunakan ibu jari dan jari telunjuk, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (9, 'gross_motor_skills',
    'Tanpa disangga oleh bantal, kursi atau dinding, dapatkah bayi duduk sendiri selama 60 detik?',
    'Duduk Mandiri',
    'Berikan kesempatan bayi duduk sendiri di alas yang aman sambil mengawasinya agar keseimbangan dan kekuatan otot tubuhnya semakin baik, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (9, 'socialization',
    'Apakah bayi dapat makan kue kering sendiri?',
    'Makan Sendiri',
    'Berikan makanan yang mudah dipegang agar bayi dapat mencoba makan sendiri sambil melatih koordinasi tangan dan kemandiriannya, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (9, 'speech_and_language',
    'Pada waktu bayi bermain sendiri dan anda diam-diam datang berdiri di belakangnya, apakah ia menengok ke belakang seperti mendengar kedatangan anda? Suara keras tidak ikut dihitung. Jawab YA hanya jika anda melihat reaksinya terhadap suara yang perlahan atau bisikan.',
    'Mengenali Suara',
    'Panggil nama bayi atau buat suara lembut dari arah yang berbeda agar ia belajar mengenali dan mencari sumber suara di sekitarnya, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (9, 'socialization',
    'Letakkan suatu mainan yang dinginkannya di luar jangkauan bayi, apakah ia mencoba mendapatkannya dengan mengulurkan lengan atau badannya?',
    'Meraih Mainan',
    'Letakkan mainan favorit sedikit di luar jangkauan agar bayi terdorong mengulurkan tangan atau menggerakkan tubuhnya untuk meraih benda tersebut, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    -- =======================================
    -- BULAN 12
    -- =======================================
    (12, 'socialization',
    'Jika anda bersembunyi di belakang sesuatu/di pojok, kemudian muncui dan menghilang secara berulang-ulang di hadapan anak, apakah ia mencari anda atau mengharapkan anda muncul kembali?',
    'Bermain Cilukba',
    'Ajak anak bermain cilukba atau bersembunyi di balik benda sambil sesekali muncul kembali agar ia belajar mengenali keberadaan orang meskipun tidak selalu terlihat, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (12, 'fine_motor_skills',
    'Letakkan pensil di telapak tangan bayi. Coba ambil pensil tersebut dengan perlahan-lahan. Sulitkah anda mendapatkan pensil itu kembali?',
    'Genggaman Kuat',
    'Berikan mainan atau benda yang aman untuk digenggam dan ajak anak mempertahankan genggamannya melalui permainan sederhana setiap hari, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (12, 'gross_motor_skills',
    'Apakah anak dapat berdiri selama 30 detik atau lebih dengan berpegangan pada kursi/meja?',
    'Berdiri Berpegangan',
    'Berikan kesempatan anak berdiri sambil berpegangan pada meja atau kursi yang kokoh agar kekuatan kaki dan keseimbangannya berkembang secara bertahap, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (12, 'speech_and_language',
    'Apakah anak dapat mengatakan 2 suku kata yang sama, misalnya: "ma-ma" , "da-da" atau "pa-pa". Jawab YA bila ia mengeluarkan salah satu suara tadi.',
    'Meniru Suara',
    'Sering ucapkan suku kata sederhana seperti "ma-ma", "pa-pa", atau "da-da" sambil mengajak anak bercakap agar ia terdorong menirukan bunyi yang didengarnya, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (12, 'gross_motor_skills',
    'Apakah anak dapat mengangkat badannya ke posisi berdiri tanpa bantuan anda?',
    'Bangun Berdiri',
    'Letakkan mainan favorit di atas permukaan yang aman agar anak terdorong menarik tubuhnya hingga berdiri tanpa banyak bantuan, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (12, 'socialization',
    'Apakah anak dapat membedakan anda dengan orang yang belum ia kenal? la akan menunjukkan sikap malu-malu atau raguragu pada saat permulaan bertemu dengan orang yang belum dikenalnya.',
    'Mengenal Orang',
    'Ajak anak bertemu anggota keluarga atau teman dalam suasana yang nyaman agar ia belajar mengenali orang yang dikenal maupun orang baru secara bertahap, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (12, 'fine_motor_skills',
    'Apakah anak dapat mengambil Benda kecil seperti kacang atau kismis, dengan meremas di antara ibu jari dan jarinya?',
    'Menjepit Benda',
    'Berikan makanan ringan berukuran kecil atau benda yang aman untuk melatih anak mengambil benda menggunakan ibu jari dan jari telunjuk, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (12, 'gross_motor_skills',
    'Apakah anak dapat duduk sendiri tanpa bantuan?',
    'Duduk Mandiri',
    'Berikan kesempatan anak duduk sendiri saat bermain di alas yang aman agar keseimbangan tubuh dan kemampuan duduk mandirinya semakin baik, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (12, 'speech_and_language',
    'Sebut 2-3 kata yang dapat ditiru oleh anak (tidak perlu kata-kata yang lengkap). Apakah ia mencoba meniru menyebutkan kata-kata tadi?',
    'Meniru Kata',
    'Ucapkan kata-kata sederhana sambil menunjuk benda atau orang di sekitar agar anak terdorong meniru ucapan dan mulai menambah kosakata pertamanya, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (12, 'fine_motor_skills',
    'Tanpa bantuan, apakah anak dapat mempertemukan dua kubus kecil yang ia pegang? Kerincingan bertangkai dan tutup panel tidak ikut dinilai.',
    'Menyatukan Kubus',
    'Sediakan dua balok atau kubus berukuran aman dan tunjukkan cara mempertemukannya agar anak belajar mengoordinasikan gerakan kedua tangannya melalui aktivitas bermain, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    -- =======================================
    -- BULAN 15
    -- =======================================
    (15, 'fine_motor_skills',
    'Tanpa bantuan, apakah anak dapat mempertemukan dua kubus kecil yang ia pegang? Kerincingan bertangkai dan tutup, panci tidak ikut dinilai',
    'Menyusun Kubus',
    'Ajak anak bermain menyusun dua balok atau kubus berukuran aman sambil memberikan contoh terlebih dahulu agar koordinasi kedua tangan dan keterampilan fine_motor_skillsnya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (15, 'gross_motor_skills',
    'Apakah anak dapat jalan sendiri atau jalan dengan berpegangan?',
    'Belajar Berjalan',
    'Berikan kesempatan anak berjalan sendiri di area yang datar dan aman sambil memberikan dukungan serta pujian agar rasa percaya diri dan keseimbangannya semakin meningkat, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (15, 'socialization',
    'Tanpa bantuan, apakah anak dapat bertepuk tangan atau melambai-lambai? Jawab TIDAK bila ia membutuhkan bantuan.',
    'Tepuk Tangan',
    'Sering ajak anak bermain tepuk tangan, melambaikan tangan, atau menirukan gerakan sederhana saat bernyanyi agar kemampuan sosial dan koordinasi geraknya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (15, 'speech_and_language',
    'Apakah anak dapat mengatakan "papa" ketika ia memanggil/melihat ayahnya, atau mengatakan "mama" jika memanggil/melihat ibunya? Jawab YA bila anak mengatakan salah satu diantaranya.',
    'Memanggil Orang',
    'Biasakan menyebut nama anggota keluarga dalam aktivitas sehari-hari dan dorong anak memanggil mereka secara langsung agar kemampuan berbicara dan berkomunikasinya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (15, 'gross_motor_skills',
    'Dapatkah anak berdiri sendiri tanpa berpegangan selama kira-kira 5 detik?',
    'Berdiri Mandiri',
    'Berikan kesempatan anak berdiri tanpa berpegangan selama beberapa detik di tempat yang aman agar kekuatan otot kaki dan keseimbangannya semakin baik, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (15, 'gross_motor_skills',
    'Dapatkan anak berdiri sendiri tanpa berpegangan selama 30 detik atau lebih?',
    'Berdiri Stabil',
    'Ajak anak bermain sambil berdiri tanpa bantuan dengan durasi yang semakin lama sesuai kemampuannya agar keseimbangan tubuh berkembang secara bertahap, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (15, 'gross_motor_skills',
    'Tanpa berpegangan atau menyentuh lantai, apakah anak dapat membungkuk untuk memungut mainan di lantai dan kemudian berdiri kembali?',
    'Membungkuk Aman',
    'Letakkan mainan di lantai dan dorong anak mengambilnya lalu kembali berdiri sendiri agar keseimbangan, koordinasi, dan kekuatan otot tubuhnya semakin terlatih, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (15, 'socialization',
    'Apakah anak dapat menunjukkan apa yang diinginkannya tanpa menangis atau merengek? Jawab YA bila ia menunjuk, menarik atau mengeluarkan suara yang menyenangkan',
    'Mengungkapkan Keinginan',
    'Biasakan memberikan pilihan sederhana kepada anak dan dorong ia menunjuk, menarik, atau mengucapkan keinginannya tanpa menangis agar kemampuan komunikasinya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (15, 'gross_motor_skills',
    'Apakah anak dapat berjalan di sepanjang ruangan tanpa jatuh atau terhuyunghuyung?',
    'Berjalan Seimbang',
    'Ajak anak berjalan di sepanjang ruangan atau halaman yang aman tanpa terburu-buru agar keseimbangan dan koordinasi langkahnya semakin baik, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (15, 'fine_motor_skills',
    'Apakah anak dapat mengambil benda kecil seperti kacang, kismis, atau potongan biskuit dengan menggunakan ibu jari?',
    'Menjepit Benda',
    'Sediakan benda kecil yang aman seperti potongan makanan sesuai usia untuk melatih anak mengambil benda menggunakan ibu jari dan jari telunjuk melalui aktivitas bermain, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    -- =======================================
    -- BULAN 18
    -- =======================================
    (18, 'socialization',
    'Tanpa bantuan, apakah anak dapat bertepuk tangan atau melambai-lambai? Jawab TIDAK bila ia membutuhkan bantuan.',
    'Tepuk Tangan',
    'Ajak anak bermain tepuk tangan, melambaikan tangan, atau mengikuti lagu dengan gerakan sederhana agar kemampuan meniru gerakan dan interaksi sosialnya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (18, 'speech_and_language',
    'Apakah anak dapat mengatakan "papa" ketika ia memanggil/melihat ayahnya, atau mengatakan "mama" jika memanggil/melihat ibunya?',
    'Memanggil Orang',
    'Biasakan mengajak anak memanggil ayah, ibu, atau anggota keluarga lainnya dalam berbagai kesempatan agar kemampuan berbicara dan berkomunikasinya semakin meningkat, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (18, 'gross_motor_skills',
    'Apakah anak dapat berdiri sendiri tanpa berpegangan selama kira-kira 5 detik?',
    'Berdiri Mandiri',
    'Berikan kesempatan anak berdiri tanpa berpegangan selama beberapa detik saat bermain agar keseimbangan dan kekuatan otot kakinya semakin baik, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (18, 'gross_motor_skills',
    'Apakah anak dapat berdiri sendiri tanpa berpegangan selama 30 detik atau lebih?',
    'Berdiri Stabil',
    'Dorong anak bermain sambil berdiri tanpa bantuan dengan durasi yang semakin lama sesuai kemampuannya agar kontrol postur dan keseimbangannya terus berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (18, 'gross_motor_skills',
    'Tanpa berpegangan atau menyentuh lantai, apakah anak dapat membungkuk untuk memungut mainan di lantai clan kemudian berdiri kembali?',
    'Membungkuk Aman',
    'Letakkan mainan di lantai dan dorong anak mengambilnya lalu kembali berdiri sendiri agar koordinasi gerak dan keseimbangan tubuhnya semakin terlatih, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (18, 'socialization',
    'Apakah anak dapat menunjukkan apa yang diinginkannya tanpa menangis atau merengek? Jawab YA bila ia menunjuk, menarik atau mengeluarkan suara yang menyenangkan.',
    'Menyampaikan Keinginan',
    'Berikan pilihan sederhana dalam kegiatan sehari-hari dan dorong anak mengungkapkan keinginannya melalui kata, menunjuk, atau gerakan tanpa menangis agar kemampuan komunikasinya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (18, 'gross_motor_skills',
    'Apakah anak dapat berjalan di sepanjang ruangan tanpa jatuh atau terhuyung-huyung ?',
    'Berjalan Stabil',
    'Ajak anak berjalan di berbagai permukaan yang aman, seperti lantai rumah atau halaman yang rata, agar keseimbangan dan koordinasi langkahnya semakin baik, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (18, 'fine_motor_skills',
    'Apakah anak anak dapat mengambil benda kecil seperti kacang, kismis, atau potongan biskuit dengan menggunakan ibu jari dan jari telunjuk?',
    'Menjepit Benda',
    'Sediakan benda kecil yang aman seperti potongan buah lunak atau makanan sesuai usia untuk melatih anak mengambil benda menggunakan ibu jari dan jari telunjuk, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (18, 'fine_motor_skills',
    'Jika anda menggelindingkan bola ke anak, apakah ia menggelindingkan /melemparkan kembali bola pada anda?',
    'Bermain Bola',
    'Luangkan waktu bermain lempar dan gelinding bola bersama anak secara bergantian agar koordinasi mata dan tangan serta kemampuan berinteraksi sosialnya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (18, 'socialization',
    'Apakah anak dapat memegang sendiri cangkir/gelas dan minum dari tempat tersebut tanpa tumpah?',
    'Minum Mandiri',
    'Berikan kesempatan anak minum menggunakan gelas sendiri saat waktu makan sambil tetap mengawasi agar kemandirian dan koordinasi gerak tangannya semakin terlatih, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    -- =======================================
    -- BULAN 21
    -- =======================================
    (21, 'gross_motor_skills',
    'Tanpa berpegangan atau menyentuh lantai, apakah anak dapat membungkuk untuk memungut mainan di lantai dan kemudian berdiri kembali?',
    'Membungkuk Mandiri',
    'Ajak anak mengambil mainan yang diletakkan di lantai lalu kembali berdiri sendiri agar keseimbangan, koordinasi gerak, dan kekuatan otot tubuhnya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (21, 'socialization',
    'Apakah anak dapat menunjukkan apa yang diinginkannya tanpa menangis atau merengek? Jawab YA bila ia menunjuk, menarik atau mengeluarkan suara yang menyenangkan.',
    'Mengungkapkan Keinginan',
    'Biasakan memberi pilihan sederhana dalam aktivitas sehari-hari dan dorong anak menyampaikan keinginannya melalui kata, menunjuk, atau gerakan tanpa menangis agar kemampuan komunikasinya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (21, 'gross_motor_skills',
    'Apakah anak dapat berjalan di sepanjang ruangan tanpa jatuh atau terhuyunghuyung ?',
    'Berjalan Mandiri',
    'Berikan kesempatan anak berjalan sendiri di lingkungan yang aman dengan pengawasan orang tua agar keseimbangan, koordinasi langkah, dan rasa percaya dirinya semakin meningkat, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (21, 'fine_motor_skills',
    'Apakah anak dapat mengambil benda kecil seperti kacang, kismis, atau potongan biskuit dengan menggunakan ibu jari clan jari telunjuk?',
    'Menjepit Benda',
    'Sediakan benda kecil yang aman sesuai usia, seperti potongan buah lunak atau balok kecil, untuk melatih anak mengambil benda menggunakan ibu jari dan jari telunjuk, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (21, 'fine_motor_skills',
    'Jika anda menggelindingkan bola ke anak, apakah ia menggelindingkan /melemparkan kembali bola pada anda?',
    'Menggelinding Bola',
    'Luangkan waktu bermain saling menggelindingkan bola dengan anak secara bergantian agar koordinasi gerak, kemampuan mengikuti permainan, dan interaksi sosialnya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (21, 'socialization',
    'Apakah anak dapat memegang sendiri cangkir/gelas clan minum dari tempat tersebut tanpa tumpah?',
    'Minum Sendiri',
    'Berikan kesempatan anak minum menggunakan gelas sendiri saat makan atau minum dengan tetap diawasi agar koordinasi tangan dan kemandiriannya semakin terlatih, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (21, 'socialization',
    'Jika anda sedang melakukan pekerjaan rumah tangga, apakah anak meniru apa yang anda lakukan?',
    'Meniru Aktivitas',
    'Libatkan anak dalam kegiatan rumah tangga sederhana, seperti menyapu dengan sapu mainan, mengelap meja, atau merapikan mainan agar kemampuan meniru dan interaksi sosialnya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (21, 'fine_motor_skills',
    'Apakah anak dapat meletakkan satu kubus di atas kubus yang lain tanpa menjatuhkan kubus itu? Kubus yang digunakan ukuran 2.5-5.0 cm',
    'Menyusun Kubus',
    'Sediakan balok atau kubus berukuran aman dan ajak anak menyusun dua hingga tiga balok secara bertahap sambil memberikan contoh agar koordinasi tangan dan keterampilan fine_motor_skillsnya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (21, 'speech_and_language',
    'Apakah anak dapat mengucapkan paling sedikit 3 kata yang mempunyai arti selain "papa" dan "mama"?.',
    'Menambah Kosakata',
    'Sering ajak anak berbicara, membaca buku bergambar, dan menyebutkan nama benda di sekitarnya agar kosakata serta kemampuan mengucapkan kata yang bermakna semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (21, 'gross_motor_skills',
    'Apakah anak dapat berjalan mundur 5 langkah atau lebih tanpa kehilangan keseimbangan? (Anda mungkin dapat melihatnya ketika anak menarik mainannya)',
    'Berjalan Mundur',
    'Ajak anak bermain mengikuti garis atau menarik mainan sambil berjalan mundur beberapa langkah di area yang aman agar keseimbangan dan koordinasi geraknya semakin baik, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    -- =======================================
    -- BULAN 24
    -- =======================================
    (24, 'socialization',
    'Jika anda sedang melakukan pekerjaan rumah tangga, apakah anak meniru apa yang anda lakukan?',
    'Meniru Kegiatan',
    'Libatkan anak dalam kegiatan rumah tangga sederhana seperti menyapu, mengelap meja, atau merapikan mainan agar kemampuan meniru dan interaksi sosialnya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (24, 'fine_motor_skills',
    'Apakah anak dapat meletakkan 1 buah kubus di atas kubus yang lain tanpa menjatuhkan kubus itu? Kubus yang digunakan ukuran 2.5-5 cm.',
    'Menyusun Balok',
    'Sediakan balok berukuran aman dan ajak anak menyusun balok satu per satu menjadi menara sambil memberikan contoh agar koordinasi tangan dan keterampilan fine_motor_skillsnya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (24, 'speech_and_language',
    'Apakah anak dapat mengucapkan paling sedikit 3 kata yang mempunyai arti selain "papa" dan "mama"?',
    'Menambah Kosakata',
    'Biasakan mengajak anak berbicara, membaca buku bergambar, dan mengenalkan nama benda di sekitarnya agar kosakata serta kemampuan mengucapkan kata yang bermakna semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (24, 'gross_motor_skills',
    'Apakah anak dapat berjalan mundur 5 langkah atau lebih tanpa kehilangan keseimbangan? (Anda mungkin dapat melihatnya ketika anak menarik mainannya).',
    'Berjalan Mundur',
    'Ajak anak bermain mengikuti garis atau menarik mainan sambil berjalan mundur beberapa langkah di area yang aman agar keseimbangan dan koordinasi geraknya semakin baik, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (24, 'fine_motor_skills',
    'Dapatkah anak melepas pakaiannya seperti: baju, rok, atau celananya? (topi dan kaos kaki tidak ikut dinilai).',
    'Melepas Pakaian',
    'Berikan kesempatan kepada anak untuk mencoba melepas baju, celana, atau rok sendiri sebelum dibantu agar kemandirian dan koordinasi gerakan tangannya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (24, 'gross_motor_skills',
    'Dapatkah anak berjalan naik tangga sendiri? Jawab YA jika ia naik tangga dengan posisi tegak atau berpegangan pada dinding atau pegangan tangga. Jawab TIDAK jika ia naik tangga dengan merangkak atau anda tidak membolehkan anak naik tangga atau anak harus berpegangan pada seseorang.',
    'Naik Tangga',
    'Dampingi anak berlatih menaiki tangga dengan berpegangan pada pegangan tangga atau dinding yang kokoh agar kekuatan kaki dan keseimbangannya berkembang secara bertahap, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (24, 'speech_and_language',
    'Tanpa bimbingan, petunjuk atau bantuan anda, dapatkah anak menunjuk dengan benar paling sedikit satu bagian badannya (rambut, mata, hidung, mulut, atau bagian badan yang lain)?',
    'Mengenal Tubuh',
    'Ajak anak bermain sambil menyebut dan menunjuk bagian-bagian tubuh, seperti mata, hidung, telinga, dan mulut, agar ia semakin mengenali anggota tubuhnya, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (24, 'socialization',
    'Dapatkah anak makan nasi sendiri tanpa banyak tumpah?',
    'Makan Mandiri',
    'Berikan kesempatan anak makan menggunakan sendok sendiri setiap kali waktu makan agar koordinasi tangan dan kemandiriannya semakin terlatih meskipun masih terdapat sedikit makanan yang tumpah, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (24, 'socialization',
    'Dapatkah anak membantu memungut mainannya sendiri atau membantu mengangkat piring jika diminta?',
    'Membantu Orang',
    'Libatkan anak dalam tugas sederhana seperti merapikan mainan atau membawa piring plastik setelah makan agar tumbuh rasa tanggung jawab dan kemandiriannya, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (24, 'gross_motor_skills',
    'Dapatkah anak menendang bola kecil (sebesar bola tenis) ke depan tanpa berpegangan pada apapun? Mendorong tidak ikut dinilai.',
    'Menendang Bola',
    'Luangkan waktu bermain menendang bola bersama anak di area yang aman agar kekuatan otot kaki, keseimbangan, dan koordinasi geraknya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    -- =======================================
    -- BULAN 30
    -- =======================================
    (30, 'socialization',
    'Dapatkah anak melepas pakaiannya seperti: baju, rok, atau celananya? (topi clan kaos kaki tidak ikut dinilai)',
    'Melepas Pakaian',
    'Berikan kesempatan anak melepas baju, celana, atau rok sendiri sebelum dibantu agar rasa percaya diri, koordinasi tangan, dan kemandiriannya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (30, 'gross_motor_skills',
    'Dapatkah anak berjalan naik tangga sendiri? Jawab YA jika ia naik tangga dengan posisi tegak atau berpegangan pada Binding atau pegangan tangga. Jawab TIDAK jika ia naik tangga dengan merangkak atau anda tidak membolehkan anak naik tangga atau anak harus berpegangan pada seseorang.',
    'Naik Tangga',
    'Dampingi anak berlatih menaiki tangga dengan langkah bergantian sambil berpegangan pada pegangan tangga hingga ia semakin percaya diri melakukannya, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (30, 'speech_and_language',
    'Tanpa bimbingan, petunjuk atau bantuan anda, dapatkah anak menunjuk dengan benar paling seclikit satu bagian badannya (rambut, mata, hidung, mulut, atau bagian badan yang lain)?',
    'Mengenal Tubuh',
    'Ajak anak bermain sambil menyebut dan menunjuk berbagai bagian tubuh menggunakan lagu, cermin, atau buku bergambar agar ia semakin mengenali tubuhnya sendiri, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (30, 'socialization',
    'Dapatkah anak makan nasi sendiri tanpa banyak tumpah?',
    'Makan Mandiri',
    'Biasakan anak makan menggunakan sendok sendiri pada setiap waktu makan agar koordinasi tangan, kemandirian, dan rasa tanggung jawab terhadap kebutuhannya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (30, 'speech_and_language',
    'Dapatkah anak membantu memungut mainannya sendiri atau membantu mengangkat piring jika diminta?',
    'Membantu Orang',
    'Libatkan anak dalam tugas sederhana seperti merapikan mainan, menyimpan buku, atau membawa peralatan makan yang ringan agar kepedulian dan kemandiriannya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (30, 'gross_motor_skills',
    'Dapatkah anak menendang bola kecil (sebesar bola tenis) ke depan tanpa berpegangan pada apapun? Mendorong tidak ikut dinilai.',
    'Menendang Bola',
    'Ajak anak bermain menendang bola ke arah sasaran atau kepada anggota keluarga secara bergantian agar koordinasi gerak, keseimbangan, dan kekuatan otot kakinya semakin baik, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (30, 'fine_motor_skills',
    'Bila diberi pensil, apakah anak mencoretcoret kertas tanpa bantuan/petunjuk?',
    'Mencoret Kertas',
    'Sediakan krayon berukuran besar dan kertas kosong agar anak bebas mencoret, membuat garis, atau bereksplorasi menggambar untuk melatih koordinasi tangan dan kontrol fine_motor_skillsnya, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (30, 'fine_motor_skills',
    'Dapatkah anak meletakkan 4 buah kubus satu persatu di atas kubus yang lain tanpa menjatuhkan kubus itu? Kubus yang digunakan ukuran 2.5-5 cm.',
    'Menyusun Menara',
    'Ajak anak menyusun empat balok atau lebih menjadi menara sambil memberikan contoh dan pujian atas setiap keberhasilannya agar koordinasi tangan dan konsentrasinya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (30, 'speech_and_language',
    'Dapatkah anak menggunakan 2 kata pada saat berbicara seperti "minta minum", "mau tidur"? "Terimakasih" dan "Dadag" tidak ikut dinilai.',
    'Menggabung Kata',
    'Biasakan berbicara menggunakan kalimat sederhana dan dorong anak menjawab dengan dua kata atau lebih dalam percakapan sehari-hari agar kemampuan berbicara dan menyusun kalimatnya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (30, 'speech_and_language',
    'Apakah anak dapat menyebut 2 diantara gambar-gambar ini tanpa bantuan?',
    'Mengenal Gambar',
    'Gunakan buku bergambar atau kartu bergambar dan ajak anak menyebutkan nama benda yang dikenalnya agar kemampuan berbahasa, daya ingat, dan pengamatan visualnya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    -- =======================================
    -- BULAN 36
    -- =======================================
    (36, 'fine_motor_skills',
    'Bila diberi pensil, apakah anak mencoret-coret kertas tanpa bantuan/petunjuk?',
    'Menggambar Bebas',
    'Sediakan krayon atau pensil warna dan ajak anak mencoret maupun menggambar bebas di atas kertas agar koordinasi tangan, kontrol gerakan, dan kreativitasnya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (36, 'fine_motor_skills',
    'Dapatkah anak meletakkan 4 buah kubus satu persatu di atas kubus yang lain tanpa menjatuhkan kubus itu? Kubus yang digunakan ukuran 2.5-5 cm.',
    'Menyusun Menara',
    'Ajak anak menyusun empat balok atau lebih menjadi menara sambil memberikan contoh dan kesempatan mencoba sendiri agar koordinasi tangan, konsentrasi, dan ketelitiannya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (36, 'speech_and_language',
    'Dapatkah anak menggunakan 2 kata pada saat berbicara seperti "minta minum"; "mau tidur"? "Terimakasih" dan "Dadag" tidak ikut dinilai.',
    'Berbicara Kalimat',
    'Biasakan mengajak anak bercakap-cakap menggunakan kalimat sederhana dan dorong ia menjawab dengan dua hingga tiga kata yang bermakna agar kemampuan berbahasanya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (36, 'speech_and_language',
    'Apakah anak dapat menyebut 2 diantara gambargambar ini tanpa bantuan?',
    'Menyebut Gambar',
    'Gunakan buku cerita atau kartu bergambar dan ajak anak menyebutkan nama benda, hewan, atau orang yang dikenalnya agar kosakata dan kemampuan berbahasanya semakin bertambah, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (36, 'gross_motor_skills',
    'Dapatkah anak melempar bola lurus ke arah perut atau dada anda dari jarak 1,5 meter?',
    'Melempar Bola',
    'Luangkan waktu bermain lempar tangkap menggunakan bola berukuran kecil secara bertahap agar koordinasi mata dan tangan serta kemampuan gross_motor_skillsnya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (36, 'speech_and_language',
    'Ikuti perintah ini dengan seksama. Jangan memberi isyarat dengan telunjuk atau mata pada saat memberikan perintah berikut ini: "Letakkan kertas ini di lantai". "Letakkan kertas ini di kursi". "Berikan kertas ini kepada ibu". Dapatkah anak melaksanakan ketiga perintah tadi?',
    'Mengikuti Perintah',
    'Berikan instruksi sederhana yang terdiri dari dua hingga tiga langkah saat bermain atau beraktivitas agar kemampuan mendengar, memahami, dan mengikuti arahan semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (36, 'fine_motor_skills',
    'Buat garis lurus ke bawah sepanjang sekurangkurangnya 2.5 cm. Suruh anak menggambar garis lain di samping garis tsb.',
    'Membuat Garis',
    'Ajak anak meniru membuat garis lurus, garis datar, atau garis tegak menggunakan pensil atau krayon agar kontrol gerakan tangan dan kesiapan menulisnya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (36, 'gross_motor_skills',
    'Letakkan selembar kertas seukuran buku di lantai. Apakah anak dapat melompati bagian lebar kertas dengan mengangkat kedua kakinya secara bersamaan tanpa didahului lari?',
    'Melompat Maju',
    'Ajak anak bermain melompati garis atau benda tipis di lantai menggunakan kedua kaki secara bersamaan agar kekuatan otot kaki, keseimbangan, dan koordinasi geraknya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (36, 'socialization',
    'Dapatkah anak mengenakan sepatunya sendiri?',
    'Memakai Sepatu',
    'Biasakan anak memakai dan melepas sepatunya sendiri sebelum keluar rumah agar rasa percaya diri dan kemandiriannya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (36, 'gross_motor_skills',
    'Dapatkah anak mengayuh sepeda roda tiga sejauh sedikitnya 3 meter?',
    'Mengayuh Sepeda',
    'Berikan kesempatan anak bermain sepeda roda tiga di area yang aman dan datar agar kekuatan kaki, keseimbangan, dan koordinasi geraknya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    -- =======================================
    -- BULAN 42
    -- =======================================
    (42, 'socialization',
    'Dapatkah anak mengenakan sepatunya sendiri?',
    'Memakai Sepatu',
    'Biasakan anak memakai dan melepas sepatunya sendiri sebelum beraktivitas di luar rumah agar rasa percaya diri, koordinasi gerakan, dan kemandiriannya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (42, 'gross_motor_skills',
    'Dapatkah anak mengayuh sepeda rods tiga sejauh sedikitnya 3 meter?',
    'Mengayuh Sepeda',
    'Berikan kesempatan anak mengayuh sepeda roda tiga di area yang datar dan aman secara rutin agar kekuatan otot kaki, keseimbangan, dan koordinasi geraknya semakin meningkat, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (42, 'socialization',
    'Setelah makan, apakah anak mencuci clan mengeringkan tangannya dengan balk sehingga anda ticlak perlu mengulanginya?',
    'Cuci Tangan',
    'Biasakan anak mencuci dan mengeringkan tangan sendiri setelah makan atau bermain agar terbentuk kemandirian sekaligus kebiasaan hidup bersih sejak dini, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (42, 'gross_motor_skills',
    'Suruh anak berdiri satu kaki tanpa berpegangan. Jika perlu tunjukkan caranya clan beri anak anda kesempatan melakukannya 3 kali. Dapatkah ia mempertahankan keseimbangan dalam waktu 2 detik atau lebih?',
    'Berdiri Satu',
    'Ajak anak bermain berdiri dengan satu kaki sambil menghitung waktu atau melalui permainan sederhana agar keseimbangan dan kekuatan otot kakinya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (42, 'gross_motor_skills',
    'Letakkan selembar kertas seukuran buku ini di lantai. Apakah anak dapat melompati panjang kertas ini dengan mengangkat kedua kakinya secara bersamaan tanpa didahului lari?',
    'Melompat Jauh',
    'Ajak anak melompati garis atau benda tipis di lantai menggunakan kedua kaki secara bersamaan melalui permainan yang menyenangkan agar koordinasi gerak dan kekuatan otot kakinya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (42, 'fine_motor_skills',
    'Jangan membantu anak clan jangan menyebut lingkaran. Suruh anak menggambar seperti contoh ini di kertas kosong yang tersedia. Dapatkah anak menggambar lingkaran?',
    'Menggambar Lingkaran',
    'Berikan contoh menggambar lingkaran kemudian biarkan anak menirunya secara bertahap menggunakan pensil atau krayon agar koordinasi tangan dan kesiapan menulisnya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (42, 'fine_motor_skills',
    'Dapatkah anak meletakkan 8 buah kubus satu persatu di atas yang lain tanpa menjatuhkan kubus tersebut? Kubus yang digunakan ukuran 2.5-5 cm.',
    'Menyusun Balok',
    'Sediakan delapan balok berukuran aman dan ajak anak menyusunnya menjadi menara sambil memberikan kesempatan mencoba sendiri agar fine_motor_skills, konsentrasi, dan ketelitiannya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (42, 'socialization',
    'Apakah anak dapat bermain petak umpet, ular naga atau permainan lain dimana ia ikut bermain clan mengikuti aturan bermain?',
    'Bermain Aturan',
    'Luangkan waktu bermain permainan sederhana yang memiliki aturan, seperti petak umpet atau permainan bergiliran, agar anak belajar mengikuti aturan, bekerja sama, dan menunggu giliran, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (42, 'socialization',
    'Dapatkah anak mengenakan celana panjang, kemeja, baju atau kaos kaki tanpa di bantu? (Tidak termasuk memasang kancing, gesper atau ikat pinggang)',
    'Berpakaian Mandiri',
    'Berikan kesempatan anak mengenakan celana, baju, atau kaus kaki sendiri sebelum dibantu agar keterampilan berpakaian dan rasa percaya dirinya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    -- =======================================
    -- BULAN 48
    -- =======================================
    (48, 'gross_motor_skills',
    'Dapatkah anak mengayuh sepeda roda tiga sejauh sedikitnya 3 meter?',
    'Mengayuh Sepeda',
    'Berikan kesempatan anak mengayuh sepeda roda tiga atau sepeda dengan roda bantu di area yang aman agar kekuatan otot kaki, koordinasi, dan keseimbangannya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (48, 'socialization',
    'Setelah makan, apakah anak mencuci dan mengeringkan tangannya dengan baik sehingga anda tidak perlu mengulanginya?',
    'Cuci Tangan',
    'Biasakan anak mencuci dan mengeringkan tangan sendiri sebelum makan, setelah dari toilet, dan setelah bermain agar terbentuk kebiasaan hidup bersih sekaligus meningkatkan kemandiriannya, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (48, 'gross_motor_skills',
    'Suruh anak berdiri satu kaki tanpa berpegangan. Jika perlu tunjukkan caranya dan beri anak anda kesempatan melakukannya 3 kali. Dapatkah ia mempertahankan keseimbangan dalam waktu 2 detik atau lebih?',
    'Berdiri Satu',
    'Ajak anak bermain berdiri dengan satu kaki sambil menghitung atau bernyanyi agar keseimbangan tubuh dan kekuatan otot kakinya semakin terlatih, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (48, 'gross_motor_skills',
    'Letakkan selembar kertas seukuran buku ini di lantai. Apakah anak dapat melompati panjang kertas ini dengan mengangkat kedua kakinya secara bersamaan tanpa didahului lari?',
    'Melompat Jauh',
    'Ajak anak melompati garis atau benda tipis menggunakan kedua kaki secara bersamaan melalui permainan yang menyenangkan agar koordinasi gerak dan kekuatan otot kakinya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (48, 'fine_motor_skills',
    'Jangan membantu anak dan jangan menyebut lingkaran. Suruh anak menggambar seperti contoh ini di kertas kosong yang tersedia. Dapatkah anak menggambar lingkaran?',
    'Menggambar Lingkaran',
    'Berikan kesempatan anak menggambar lingkaran menggunakan pensil atau krayon tanpa banyak bantuan agar koordinasi mata dan tangan serta kesiapan menulisnya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (48, 'fine_motor_skills',
    'Dapatkah anak meletakkan 8 buah kubus satu persatu di atas yang lain tanpa menjatuhkan kubus tersebut? Kubus yang digunakan ukuran 2.5-5 cm.',
    'Menyusun Balok',
    'Sediakan delapan balok berukuran aman dan ajak anak menyusunnya menjadi menara sambil memberi kesempatan mencoba sendiri agar fine_motor_skills, ketelitian, dan konsentrasinya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (48, 'socialization',
    'Apakah anak dapat bermain petak umpet, ular naga atau permainan lain dimana ia ikut bermain dan mengikuti aturan bermain?',
    'Bermain Bersama',
    'Luangkan waktu bermain permainan yang memiliki aturan sederhana bersama keluarga atau teman sebaya agar anak belajar bekerja sama, menunggu giliran, dan mengikuti aturan permainan, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (48, 'socialization',
    'Dapatkah anak mengenakan celana panjang, kemeja, baju atau kaos kaki tanpa di bantu? (Tidak termasuk memasang kancing, gesper atau ikat pinggang)',
    'Berpakaian Mandiri',
    'Biasakan anak mengenakan baju, celana, dan kaus kaki sendiri setiap hari agar keterampilan berpakaian dan rasa percaya dirinya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (48, 'speech_and_language',
    'Dapatkah anak menyebutkan nama lengkapnya tanpa dibantu? Jawab TIDAK jika ia hanya menyebutkan sebagian namanya atau ucapannya sulit dimengerti.',
    'Menyebut Nama',
    'Biasakan memperkenalkan nama lengkap anak dalam percakapan sehari-hari dan dorong anak menyebutkannya sendiri saat diminta agar kemampuan berbahasa dan rasa percaya dirinya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    -- =======================================
    -- BULAN 54
    -- =======================================
    (54, 'fine_motor_skills',
    'Dapatkah anak meletakkan 8 buah kubus satu persatu di atas yang lain tanpa menjatuhkan kubus tersebut? Kubus yang digunakan ukuran 2-5-5',
    'Menyusun Balok',
    'Sediakan delapan balok berukuran aman dan ajak anak menyusunnya menjadi menara sambil memberi kesempatan mencoba sendiri agar koordinasi tangan, ketelitian, dan konsentrasinya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (54, 'socialization',
    'Apakah anak dapat bermain petak umpet, ular naga atau permainan lain dimana ia ikut bermain dan mengikuti aturan bermain?',
    'Bermain Aturan',
    'Luangkan waktu bermain permainan yang memiliki aturan sederhana bersama keluarga atau teman sebaya agar anak belajar bekerja sama, menunggu giliran, dan mengikuti aturan permainan, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (54, 'socialization',
    'Dapatkah anak mengenakan celana panjang, kemeja, baju atau kaos kaki tanpa di bantu? (Tidak termasuk memasang kancing, gesper atau ikat pinggang)',
    'Berpakaian Mandiri',
    'Biasakan anak mengenakan baju, celana, dan kaus kaki sendiri setiap hari agar keterampilan berpakaian dan rasa percaya dirinya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (54, 'speech_and_language',
    'Dapatkah anak menyebutkan nama lengkapnya tanpa dibantu? Jawab TIDAK jika ia hanya menyebut sebagian namanya atau ucapannya sulit dimengerti.',
    'Menyebut Nama',
    'Biasakan memperkenalkan nama lengkap anak dalam percakapan sehari-hari dan dorong anak menyebutkannya sendiri saat diminta agar kemampuan berbahasa dan rasa percaya dirinya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (54, 'speech_and_language',
    'Isi titik-titik di bawah ini dengan jawaban anak. Jangan membantu kecuali mengulangi pertanyaan. "Apa yang kamu lakukan jika kamu kedinginan?" "Apa yang kamu lakukan jika kamu lapar?" "Apa yang kamu lakukan jika kamu lelah?" Jawab YA biia anak merjawab ke 3 pertanyaan tadi dengan benar, bukan dengan gerakan atau isyarat. Jika kedinginan, jawaban yang benar adalah "menggigil", "pakai mantel'' atau "masuk kedalam rumah''. Jika lapar, jawaban yang benar adalah "makan" Jika lelah, jawaban yang benar adalah "mengantuk", "tidur", "berbaring/tidur-tiduran", "istirahat" atau "diam sejenak"',
    'Menjawab Situasi',
    'Ajak anak berdiskusi tentang apa yang dilakukan saat lapar, lelah, atau kedinginan agar kemampuan memahami situasi dan mengungkapkan jawaban dengan kata-kata semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (54, 'socialization',
    'Apakah anak dapat mengancingkan bajunya atau pakaian boneka?',
    'Mengancing Baju',
    'Berikan kesempatan anak berlatih membuka dan mengancingkan baju atau pakaian boneka melalui permainan sehari-hari agar koordinasi jari dan kemandiriannya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (54, 'gross_motor_skills',
    'Suruh anak berdiri satu kaki tanpa berpegangan. Jika perlu tunjukkan caranya dan beri anak ands kesempatan melakukannya 3 kali. Dapatkah ia mempertahankan keseimbangan dalam waktu 6 detik atau lebih?',
    'Berdiri Satu',
    'Ajak anak bermain berdiri dengan satu kaki sambil menghitung waktu atau bernyanyi agar keseimbangan tubuh dan kekuatan otot kakinya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (54, 'fine_motor_skills',
    'Jangan mengoreksi/membantu anak. Jangan menyebut kata "lebih panjang". Perlihatkan gambar kedua garis ini pada anak. Tanyakan: "Mana garis yang lebih panjang?" Minta anak menunjuk garis yang lebih panjang. Setelah anak menunjuk, putar lembar ini dan ulangi pertanyaan tersebut. Setelah anak menunjuk, putar lembar ini lagi dan ulangi pertanyaan tadi. Apakah anak dapat menunjuk garis yang lebih panjang sebanyak 3 kali dengan benar?',
    'Membanding Garis',
    'Ajak anak bermain membandingkan benda yang lebih panjang dan lebih pendek menggunakan pensil, tali, atau balok agar kemampuan mengamati dan memahami perbedaan ukuran semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (54, 'fine_motor_skills',
    'Jangan membantu anak dan jangan memberitahu nama gambar ini, suruh anak menggambar seperti contoh ini di kertas kosong yang tersedia. Berikan 3 kali kesempatan. Apakah anak dapat menggambar seperti contoh ini?',
    'Meniru Gambar',
    'Berikan contoh gambar sederhana lalu biarkan anak menirunya menggunakan pensil atau krayon agar koordinasi mata dan tangan serta kemampuan fine_motor_skillsnya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (54, 'speech_and_language',
    'Ikuti perintah ini dengan seksama. Jangan memberi isyarat dengan telunjuk atau mats pads saat memberikan perintah berikut ini: "Letakkan kertas ini di atas lantai". "Letakkan kertas ini di bawah kursi". "Letakkan kertas ini di depan kamu" "Letakkan kertas ini di belakang kamu" Jawab YA hanya jika anak mengerti arti "di atas", "di bawah", "di depan" dan "di belakang"',
    'Memahami Posisi',
    'Gunakan permainan menyimpan benda di atas, bawah, depan, dan belakang agar anak belajar memahami dan mengikuti petunjuk tentang posisi benda dalam ruang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    -- =======================================
    -- BULAN 60
    -- =======================================
    (60, 'speech_and_language',
    'Isi titik-titik di bawah ini dengan jawaban anak. Jangan membantu kecuali mengulangi pertanyaan. "Ара yang kamu lakukan jika kamu kedinginan?" "Ара yang kamu lakukan jika kamu lapar?" "Apa yang kamu lakukan jika kamu lelah?" Jawab YA biia anak merjawab ke 3 pertanyaan tadi dengan benar, bukan dengan gerakan atau isyarat. Jika kedinginan, jawaban yang benar adalah "menggigil", "pakai mantel'' atau "masuk kedalam rumah''. Jika lapar, jawaban yang benar adalah "makan" Jika lelah, jawaban yang benar adalah "mengantuk", "tidur", "berbaring/tidur-tiduran", "istirahat" atau "diam sejenak"',
    'Menjawab Situasi',
    'Ajak anak berdiskusi tentang apa yang dilakukan saat lapar, lelah, atau kedinginan dalam kegiatan sehari-hari agar kemampuan berpikir, memahami situasi, dan mengungkapkan jawaban dengan kata-kata semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (60, 'socialization',
    'Apakah anak dapat mengancingkan bajunya atau pakaian boneka?',
    'Mengancing Baju',
    'Biasakan anak membuka dan mengancingkan bajunya sendiri saat berpakaian agar koordinasi jari, ketelitian, dan kemandiriannya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (60, 'gross_motor_skills',
    'Suruh anak berdiri satu kaki tanpa berpegangan. Jika perlu tunjukkan caranya dan beri anak ands kesempatan melakukannya 3 kali. Dapatkah ia mempertahankan keseimbangan dalam waktu 6 detik atau lebih?',
    'Berdiri Satu',
    'Ajak anak bermain berdiri dengan satu kaki sambil menghitung hingga enam detik atau melalui permainan keseimbangan lainnya agar kekuatan otot dan keseimbangan tubuhnya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (60, 'fine_motor_skills',
    'Jangan mengoreksi/membantu anak. Jangan menyebut kata "lebih panjang". Perlihatkan gambar kedua garis ini pada anak. Tanyakan: "Mana garis yang lebih panjang?" Minta anak menunjuk garis yang lebih panjang. Setelah anak menunjuk, putar lembar ini dan ulangi pertanyaan tersebut. Setelah anak menunjuk, putar lembar ini lagi dan ulangi pertanyaan tadi. Apakah anak dapat menunjuk garis yang lebih panjang sebanyak 3 kali dengan benar?',
    'Membanding Garis',
    'Gunakan permainan membandingkan panjang berbagai benda di sekitar rumah agar anak belajar mengenali konsep lebih panjang dan lebih pendek melalui pengalaman yang menyenangkan, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (60, 'fine_motor_skills',
    'Jangan membantu anak dan jangan memberitahu nama gambar ini, suruh anak menggambar seperti contoh ini di kertas kosong yang tersedia. Berikan 3 kali kesempatan. Apakah anak dapat menggambar seperti contoh ini?',
    'Meniru Gambar',
    'Berikan contoh gambar sederhana kemudian biarkan anak menirunya tanpa banyak bantuan agar koordinasi mata dan tangan serta kemampuan fine_motor_skillsnya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (60, 'speech_and_language',
    'Ikuti perintah ini dengan seksama. Jangan memberi isyarat dengan telunjuk atau mats pads saat memberikan perintah berikut ini: "Letakkan kertas ini di atas lantai". "Letakkan kertas ini di bawah kursi". "Letakkan kertas ini di depan kamu" "Letakkan kertas ini di belakang kamu" Jawab YA hanya jika anak mengerti arti "di atas", "di bawah", "di depan" dan "di belakang"',
    'Memahami Posisi',
    'Ajak anak bermain mengikuti petunjuk seperti meletakkan benda di atas, di bawah, di depan, atau di belakang agar kemampuan memahami konsep ruang dan mengikuti instruksi semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (60, 'socialization',
    'Apakah anak bereaksi dengan tenang dan tidak rewel (tanpa menangis atau menggelayut pada anda) pada saat anda meninggalkannya?',
    'Berpisah Tenang',
    'Biasakan anak mengikuti kegiatan singkat bersama keluarga atau teman tanpa didampingi terus-menerus agar rasa percaya diri, kemandirian, dan kemampuan beradaptasinya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (60, 'speech_and_language',
    'Jangan menunjuk, membantu atau membetulkan, katakan pada anak: "Tunjukkan segi empat merah" "Tunjukkan segi empat kuning" ''Tunjukkan segi empat biru" "Tunjukkan segi empat hijau" Dapatkah anak menunjuk keempat warna itu dengan benar?',
    'Mengenal Warna',
    'Gunakan permainan mengelompokkan atau mencari benda berdasarkan warna merah, kuning, hijau, dan biru agar kemampuan mengenali warna serta memahami instruksi semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (60, 'gross_motor_skills',
    'Suruh anak melompat dengan satu kaki beberapa kali tanpa berpegangan (lompatan dengan dua kaki tidak ikut dinilai). Apakah ia dapat melompat 2-3 kali dengan satu kakinya',
    'Melompat Satu',
    'Ajak anak bermain melompat menggunakan satu kaki secara bergantian di area yang aman agar keseimbangan, koordinasidan kekuatan otot berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.'),

    (60, 'socialization',
    'Dapatkah anak sepenuhnya berpakaian sendiri tanpa bantuan?',
    'Berpakaian Sendiri',
    'Berikan kesempatan anak mengenakan seluruh pakaiannya sendiri sejak memilih hingga memakai pakaian agar kemandirian, koordinasi gerak, dan rasa percaya dirinya semakin berkembang, kemudian konsultasikan dengan tenaga kesehatan apabila kemampuan ini belum berkembang sesuai usianya.')

    ) AS t(month_target, developmental_domain, question_text, title, action_text)
),

inserted_questions AS (
    INSERT INTO assessment_kpsp_questions (month_target, developmental_domain, question_text)
    SELECT month_target, CAST(developmental_domain AS developmental_domain), question_text
    FROM data
    RETURNING id, month_target, question_text
)
INSERT INTO recommended_actions (assessment_kpsp_question_id, title, action_text)
SELECT iq.id, d.title, d.action_text
FROM data d
JOIN inserted_questions iq
  ON d.month_target = iq.month_target
 AND d.question_text = iq.question_text;
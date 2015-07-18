package parser_test

const (
	goodLog = `
10.10.10.15 - - [28/Jun/2015:07:50:55 -0700] "GET /wpad.dat HTTP/1.1" 304 0 "-" "WinHttp-Autoproxy-Service/5.1"
103.28.226.211 - - [30/Jun/2015:19:38:35 -0700] "GET / HTTP/1.1" 200 396 "-" "Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko"
74.202.182.46 - - [30/Jun/2015:21:55:39 -0700] "GET /wp-login.php HTTP/1.1" 404 208 "-" "Mozilla/5.0 (Windows NT 5.1; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/24.0.1290.1 Safari/537.13"
74.202.182.46 - - [30/Jun/2015:21:55:39 -0700] "GET /administrator/ HTTP/1.1" 404 208 "-" "Mozilla/5.0 (Windows NT 5.1; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/24.0.1290.1 Safari/537.13"
74.202.182.46 - - [30/Jun/2015:21:55:39 -0700] "GET /admin.php HTTP/1.1" 404 208 "-" "Mozilla/5.0 (Windows NT 5.1; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/24.0.1290.1 Safari/537.13"
119.147.146.189 - - [30/Jun/2015:22:41:24 -0700] "GET / HTTP/1.1" 200 396 "-" "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95  Safari/537.36"
38.111.147.83 - - [30/Jun/2015:23:44:41 -0700] "GET / HTTP/1.0" 200 612 "-" "TurnitinBot/3.0 (http://www.turnitin.com/robot/crawlerinfo.html)"
219.243.6.213 - - [01/Jul/2015:00:06:49 -0700] "GET / HTTP/1.1" 200 396 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:38.0) Gecko/20100101 Firefox/38.0"
121.52.223.239 - - [01/Jul/2015:01:22:11 -0700] "GET / HTTP/1.1" 200 396 "-" "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.124 Safari/537.36"
173.193.46.236 - - [01/Jul/2015:01:34:43 -0700] "GET /check_user.cgi HTTP/1.1" 404 579 "-" "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)"
95.215.9.222 - - [01/Jul/2015:03:10:44 -0700] "GET / HTTP/1.0" 200 612 "-" "netscan.lekus.su"
58.60.231.83 - - [01/Jul/2015:03:32:45 -0700] "GET / HTTP/1.0" 200 612 "-" "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.124 Safari/537.36"
5.39.216.121 - - [01/Jul/2015:03:42:42 -0700] "GET / HTTP/1.0" 200 612 "-" "Internet-wide-scan-to-be-removed-from-this-list-email-info-at-binaryedge.io"
86.57.158.15 - - [01/Jul/2015:04:35:30 -0700] "GET /filter/tips HTTP/1.1" 404 177 "-" "Opera/9.64 (Windows NT 5.1; U; en) Presto/2.1.1"
188.115.165.152 - - [01/Jul/2015:05:15:11 -0700] "GET /wp-login.php HTTP/1.1" 404 208 "-" "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/23.0.1271.17 Safari/537.11"
122.96.17.131 - - [01/Jul/2015:05:56:54 -0700] "GET / HTTP/1.1" 200 396 "-" "Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:38.0) Gecko/20100101 Firefox/38.0"
86.57.158.15 - - [01/Jul/2015:06:25:07 -0700] "GET /filter/tips HTTP/1.1" 404 177 "-" "Opera/9.64 (Windows NT 5.1; U; en) Presto/2.1.1"
5.165.76.37 - - [01/Jul/2015:11:34:12 -0700] "GET /wp-admin/admin-ajax.php?action=revslider_show_image&img=../wp-config.php HTTP/1.0" 404 579 "-" "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; MRA 4.3 (build 01218); .NET CLR 1.1.4322; Hotbar 4.6.1)"
104.167.111.208 - - [01/Jul/2015:13:07:07 -0700] "GET / HTTP/1.1" 200 396 "http://hvd-store.com/" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.120 Safari/537.36"
199.217.117.90 - - [01/Jul/2015:15:07:25 -0700] "GET / HTTP/1.1" 200 396 "-" "Mozilla/5.0 (Windows NT 5.1; rv:7.0.1) Gecko/20100101 Firefox/7.0.1"
146.185.239.100 - - [01/Jul/2015:19:54:33 -0700] "GET http://24x7-allrequestsallowed.com/?PHPSESSID=a45ad32b00143PRTJWTIYELFQ%40 HTTP/1.1" 200 612 "-" "-"
177.205.31.209 - - [01/Jul/2015:20:23:30 -0700] "GET / HTTP/1.1" 200 396 "http://buttons-for-website.com" "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.143 Safari/537.36"
220.181.132.197 - - [01/Jul/2015:23:25:47 -0700] "GET / HTTP/1.1" 200 396 "-" "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.48 Safari/537.36"
5.9.44.13 - - [02/Jul/2015:03:11:21 -0700] "GET / HTTP/1.1" 200 396 "-" "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.124 Safari/537.36"
207.90.2.11 - - [02/Jul/2015:07:35:19 -0700] "GET / HTTP/1.0" 200 612 "-" "=Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/534.16 (KHTML, like Gecko) Chrome/10.0.648.204 Safari/534.16"
89.212.86.109 - - [02/Jul/2015:10:16:55 -0700] "GET / HTTP/1.1" 200 396 "http://videos-for-your-business.com" "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.143 Safari/537.36"
95.173.183.52 - - [02/Jul/2015:13:29:24 -0700] "OPTIONS / HTTP/1.0" 405 181 "-" "-"
37.187.129.166 - - [02/Jul/2015:19:10:04 -0700] "GET / HTTP/1.1" 200 396 "http://burger-imperia.com/" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.120 Safari/537.36"
211.103.220.197 - - [02/Jul/2015:19:46:15 -0700] "GET / HTTP/1.1" 200 396 "-" "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.130 Safari/537.36"
113.247.42.246 - - [02/Jul/2015:20:01:12 -0700] "GET / HTTP/1.1" 200 396 "-" "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 Safari/537.36"
190.7.88.158 - - [03/Jul/2015:01:49:47 -0700] "GET /tmUnblock.cgi HTTP/1.1" 404 177 "-" "-"
190.7.88.158 - - [03/Jul/2015:01:49:50 -0700] "GET / HTTP/1.1" 200 612 "-" "-"
94.142.245.231 - - [03/Jul/2015:03:08:05 -0700] "GET / HTTP/1.1" 200 396 "http://hvd-store.com/" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.120 Safari/537.36"
94.102.49.169 - - [03/Jul/2015:04:45:08 -0700] "GET / HTTP/1.1" 200 396 "-" "python-requests/2.7.0 CPython/2.7.6 Linux/3.13.0-24-generic"
95.215.9.222 - - [03/Jul/2015:06:00:10 -0700] "GET / HTTP/1.0" 200 612 "-" "netscan.lekus.su"
71.170.222.94 - - [03/Jul/2015:09:25:33 -0700] "GET /admin/config.php HTTP/1.1" 404 579 "-" "Mozilla/4.0 (compatible; MSIE 6.0; Windows 98)"
71.170.222.94 - - [03/Jul/2015:09:25:33 -0700] "GET /freepbx/admin/config.php HTTP/1.1" 404 579 "-" "Mozilla/4.0 (compatible; MSIE 6.0; Windows 98)"
66.249.85.208 - - [03/Jul/2015:11:44:02 -0700] "GET / HTTP/1.1" 200 396 "-" "Mozilla/5.0 (en-us) AppleWebKit/537.36 (KHTML, like Gecko; Google Domains) Chrome/27.0.1453 Safari/537.36"
198.154.63.131 - - [03/Jul/2015:12:30:37 -0700] "POST //%63%67%69%2D%62%69%6E/%70%68%70?%2D%64+%61%6C%6C%6F%77%5F%75%72%6C%5F%69%6E%63%6C%75%64%65%3D%6F%6E+%2D%64+%73%61%66%65%5F%6D%6F%64%65%3D%6F%66%66+%2D%64+%73%75%68%6F%73%69%6E%2E%73%69%6D%75%6C%61%74%69%6F%6E%3D%6F%6E+%2D%64+%64%69%73%61%62%6C%65%5F%66%75%6E%63%74%69%6F%6E%73%3D%22%22+%2D%64+%6F%70%65%6E%5F%62%61%73%65%64%69%72%3D%6E%6F%6E%65+%2D%64+%61%75%74%6F%5F%70%72%65%70%65%6E%64%5F%66%69%6C%65%3D%70%68%70%3A%2F%2F%69%6E%70%75%74+%2D%64+%63%67%69%2E%66%6F%72%63%65%5F%72%65%64%69%72%65%63%74%3D%30+%2D%64+%63%67%69%2E%72%65%64%69%72%65%63%74%5F%73%74%61%74%75%73%5F%65%6E%76%3D%30+%2D%64+%61%75%74%6F%5F%70%72%65%70%65%6E%64%5F%66%69%6C%65%3D%70%68%70%3A%2F%2F%69%6E%70%75%74+%2D%6E HTTP/1.1" 404 177 "-" "-"
212.16.104.33 - - [03/Jul/2015:20:01:57 -0700] "GET / HTTP/1.1" 200 396 "http://hvd-store.com/" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.120 Safari/537.36"
104.131.211.209 - - [03/Jul/2015:21:32:26 -0700] "GET / HTTP/1.1" 200 396 "-" "Mozilla/5.0 (compatible; spbot/4.4.2; +http://OpenLinkProfiler.org/bot )"
76.74.97.4 - - [04/Jul/2015:01:04:28 -0700] "GET / HTTP/1.1" 200 612 "-" "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)"
128.199.69.78 - - [04/Jul/2015:01:40:52 -0700] "GET / HTTP/1.1" 200 396 "-" "CRAZYWEBCRAWLER 0.9.6, http://www.crazywebcrawler.com"
118.192.3.4 - - [04/Jul/2015:05:11:33 -0700] "GET / HTTP/1.1" 200 396 "-" "python-requests/2.4.3 CPython/2.7.3 Linux/3.2.0-23-generic"
141.212.122.170 - - [04/Jul/2015:11:38:16 -0700] "CONNECT google.com:443 HTTP/1.1" 400 181 "-" "-"
`

	mixedLog = `
220.181.132.197 - - [01/Jul/2015:23:25:47 -0700] "GET / HTTP/1.1" 200 396 "-" "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.48 Safari/537.36"
5.9.44.13 - - [02/Jul/2015:03:11:21 -0700] "GET / HTTP/1.1" 200 396 "-" "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.124 Safari/537.36"
207.90.2.11 - - [02/Jul/2015:07:35:19 -0700] "GET / HTTP/1.0" 200 612 "-" "=Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/534.16 (KHTML, like Gecko) Chrome/10.0.648.204 Safari/534.16"
89.212.86.109 - - [02/Jul/2015:10:16:55 -0700] "GET / HTTP/1.1" 200 396 "http://videos-for-your-business.com" "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.143 Safari/537.36"
asdf
==> another.file <==
95.173.183.52 - - [02/Jul/2015:13:29:24 -0700] "OPTIONS / HTTP/1.0" 405 181 "-" "-"
37.187.129.166 - - [02/Jul/2015:19:10:04 -0700] "GET / HTTP/1.1" 200 396 "http://burger-imperia.com/" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.120 Safari/537.36"
211.103.220.197 - - [02/Jul/2015:19:46:15 -0700] "GET / HTTP/1.1" 200 396 "-" "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.130 Safari/537.36"
113.247.42.246 - - [02/Jul/2015:20:01:12 -0700] "GET / HTTP/1.1" 200 396 "-" "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 Safari/537.36"
`
)
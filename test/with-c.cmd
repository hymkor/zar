rem **** create test with -C 
setlocal
set PROMPT=$g
set "ZAR=%~dp0..\zar.exe"

pushd "%TEMP%"

rem **** create ****
mkdir zar-test
echo ahaha> zar-test\ahaha
echo ihihi> zar-test\ihihi
mkdir zar-test\uhahaha
echo gogogo> zar-test\uhahaha\gogogo
md5tree zar-test > zar-test-original.md5
%ZAR% -C zar-test -cvf ../zar-test.zip .
rmdir /s /q zar-test

rem **** extract ****
mkdir zar-test
%ZAR% -xvf zar-test.zip -Czar-test
md5tree zar-test > zar-test-extract.md5
rmdir /s /q zar-test

rem **** diff original and extract files ***
busybox diff zar-test-original.md5 zar-test-extract.md5 && echo Files are same

del zar-test-original.md5 zar-test-extract.md5
popd
endlocal
exit /b

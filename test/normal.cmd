rem **** create/extract test without -C and --remove-files ****
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
%ZAR% cvf zar-test.zip zar-test
rmdir /s /q zar-test

rem **** extract ****
%ZAR% xvf zar-test.zip
md5tree zar-test > zar-test-extract.md5
rmdir /s /q zar-test

rem **** diff original and extract files ***
busybox diff zar-test-original.md5 zar-test-extract.md5 && echo Files are same

del zar-test-original.md5 zar-test-extract.md5
popd
endlocal
exit /b

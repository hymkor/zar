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
%ZAR% --remove-files -cvf zar-test.zip -C zar-test .
rem *** If --remove-files does not delete all files, rmdir will fail ***
rmdir zar-test
del zar-test.zip

popd
endlocal
exit /b

t=$(date +%s)
while true
do
  curl https://bluegreen.mxsyx.site
  tt=$(date +%s)
  t=$tt
  echo -e '\r'
  echo `expr $tt - $t`
  echo -e "\n"
done

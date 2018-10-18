#!/bin/bash
# vim: set ts=2 sw=2 expandtab:

langs="en-us ko-kr"

for lang in $langs; do
  ls locales/*.$lang.yaml > /dev/null 2>&1 || continue
  echo -e "\n\nLANG: $lang"

  find templates -type f |xargs cat \
  |sed 's/t("/\nt("/g' \
  |grep 't("' \
  |sed 's/.*t("\([\._A-Za-z0-9]*\)").*/- id: \1Xtranslation: \1/' \
  |sed 's/\^ \([^\^]*:[^\^]*\)/\^ "\1"/g;s/\^/:/;s/\^/X/;s/\^/:/' \
  |while read line; do
    id=`echo $line|cut -dX -f1`
    grep -q "^$id\$" locales/*.$lang.yaml || echo "$line"
  done | sort -u | sed 's/Xtranslation:/\n  translation:/'

  #grep 't(c.*"\w*\(\.\w*\)*\w*"' `find actions -type f` \
  grep ' t(c, ".*' `find actions -type f` \
  |grep -v "\.html" \
  |sed 's/.*"\([^"]*\)").*/- id: \1Xtranslation: \1/' \
  |sed 's/\^ \([^\^]*:[^\^]*\)/\^ "\1"/g;s/\^/:/;s/\^/X/;s/\^/:/' \
  |while read line; do
    id=`echo $line|cut -dX -f1`
    grep -q "^$id\$" locales/*.$lang.yaml || echo "$line"
  done | sort -u | sed 's/Xtranslation:/\n  translation:/'

  echo ""
  echo "## From error strings:"

  grep 'errors.New(".*' `find actions -type f` \
  |grep -v "\.html" \
  |sed 's/.*errors.New("\([^"]*\)").*/- id^ \1^translation^ \1/' \
  |sed 's/\^ \([^\^]*:[^\^]*\)/\^ "\1"/g;s/\^/:/;s/\^/X/;s/\^/:/' \
  |while read line; do
    id=`echo $line|cut -dX -f1`
    grep -q "^$id\$" locales/*.$lang.yaml || echo "$line"
  done | sort -u | sed 's/Xtranslation:/\n  translation:/'

  grep 'errors.New(".*' `find models -type f` \
  |grep -v "\.html" \
  |sed 's/.*errors.New("\([^"]*\)").*/- id^ \1^translation^ \1/' \
  |sed 's/\^ \([^\^]*:[^\^]*\)/\^ "\1"/g;s/\^/:/;s/\^/X/;s/\^/:/' \
  |while read line; do
    id=`echo $line|cut -dX -f1`
    grep -q "^$id\$" locales/*.$lang.yaml || echo "$line"
  done | sort -u | sed 's/Xtranslation:/\n  translation:/'
done

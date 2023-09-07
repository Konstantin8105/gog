#!/bin/bash
rm -f *.dxf
rm -f *.out
rm -f *.test

echo -e "# gog"
echo -e "golang geometry library between point and segments"
echo -e "\`\`\`\n"

go doc -all .

echo -e "\`\`\`"

rm -f *.dxf
rm -f *.out
rm -f *.test

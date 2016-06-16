
rm -fr fixtures

mkdir -p fixtures/create/test1
cat <<EOT >> fixtures/create/test1/some.file
some content
EOT

mkdir -p fixtures/create/test2
mkdir -p fixtures/create/test2/sub
cat <<EOT >> fixtures/create/test2/some.file
some content
EOT
cat <<EOT >> fixtures/create/test2/sub/else.file
some other
EOT

mkdir -p fixtures/extract/test1
cat <<EOT >> fixtures/extract/test1/some.file
some content
EOT

tar -cvf fixtures/extract/test1.fedora23.tar fixtures/extract/test1

tar -cvf fixtures/extract/test2.fedora23.tar fixtures/extract/test1 -C fixtures/extract/test1

tar -zcvf fixtures/extract/test1.fedora23.tar.gz fixtures/extract/test1
cp fixtures/extract/test1.fedora23.tar.gz fixtures/extract/test1.fedora23.tgz

tar -zcvf fixtures/extract/test2.fedora23.tar.gz fixtures/extract/test1 -C fixtures/extract/test1
cp fixtures/extract/test2.fedora23.tar.gz fixtures/extract/test2.fedora23.tgz

cd fixtures/extract/test1
zip ../test1.fedora23.zip some.file
cd ../../../

zip -r fixtures/extract/test2.fedora23.zip fixtures/extract/test1

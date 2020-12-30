echo "Downloading dependencies"
go mod download
echo "Initializing fluffy"
go install
echo "done"
fluffy start &
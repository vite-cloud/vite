if [ -f ../stubs_test.go ]; then
  rm ../stubs_test.go
fi


printf "package config" >> ../stubs_test.go
printf "\n\n" >> ../stubs_test.go

# get last commit
commit=$(git rev-parse HEAD | tr -d '\n')

for dir in *; do
  if ! [ -d "$dir" ]; then
    continue
  fi

  # transform dir name in PascalCase
  name=$(echo "$dir" | sed -r 's/(^|[_-])([a-z])/\U\2/g')

  {
    printf "var %s = Locator{\n" "$name";
    printf "	Provider:   GitHubProvider{},\n";
    printf "	Path:       \"core/domain/config/testdata\",\n";
    printf "	UseHTTPS:   false,\n";
    printf "	Branch:     \"main\",\n";
    printf "	Commit:     \"%s\",\n" "$commit";
    printf "	Repository: \"vite-cloud/vite\",\n";
    printf "}\n\n"


  } >> ../stubs_test.go
done
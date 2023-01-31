
osascript -e 'tell app "Terminal" to do script "cd /Users/matheusfranco/Documents/Academica/Computacao/Tese/ssv-spec-AleaBFT/main;go run main.go launch -n 1"'
for i in {1..5}
do
  go run main.go propose -n 1
done
osascript -e 'tell app "Terminal" to do script "cd /Users/matheusfranco/Documents/Academica/Computacao/Tese/ssv-spec-AleaBFT/main;go run main.go launch -n 2"'
for i in {1..5}
do
  go run main.go propose -n 2
done
osascript -e 'tell app "Terminal" to do script "cd /Users/matheusfranco/Documents/Academica/Computacao/Tese/ssv-spec-AleaBFT/main;go run main.go launch -n 3"'
for i in {1..5}
do
  go run main.go propose -n 3
done
osascript -e 'tell app "Terminal" to do script "cd /Users/matheusfranco/Documents/Academica/Computacao/Tese/ssv-spec-AleaBFT/main;go run main.go launch -n 4"'
for i in {1..5}
do
  go run main.go propose -n 4
done
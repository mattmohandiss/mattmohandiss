# COPY STYLE
cp style.css out/style.css

# BUILD HTML
pandoc resume.md -c style.css -s --template template.html -o out/resume.html

# BUILD PDF
pandoc resume.md -c style.css -s --template template.html -t html -o out/resume.pdf

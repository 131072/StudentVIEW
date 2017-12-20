wget --directory-prefix=src -i jsDepends.txt -N
cd public
GOOS=linux gopherjs build -o ../src/studentv.js -m
cd ..
java -jar closure-compiler.jar --compilation_level=SIMPLE -W QUIET --source_map_input=src/studentv.js\|src/studentv.js.map --create_source_map public/StudentVIEW.js.map --js_output_file=public/StudentVIEW.js `find ./src/*.js`
echo '//# sourceMappingURL=StudentVIEW.js.map' >> public/StudentVIEW.js

'use strict';

define(function(require){
    var st = require('st');
    console.log(st);
    var tagger = st.NewTagger;
    alert(tagger.process("package abc\ntype a struct {\nHeight int\nWidth int\n}"))
});
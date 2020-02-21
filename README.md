# jbovlaste
a lojban dictionary parser

This is a command line program that parses the lojban language dictionary

that anyone can download as an "xml-export.html" file from this URI :

http://jbovlaste.lojban.org/export/xml-export.html?lang=en

Once compiled you can run this program from a commande line terminal
bu using this syntax :
jbovlaste [-t type][-c clue][-u user][-o][-w][-n] arg

where 'type', 'clue', 'user' and 'arg' are regular expressions
for low case strings, 'arg' being a requested argument

[-t] limit the results to entries matching a regexp in words types
[-c] limit the results to entries matching a regexp in definitions or notes
[-u] limit the results to entries matching a regexp in user names
[-o] limit the results to official data (overwrite user flag)
[-w] limit the results to wild data tagged as unofficial
[-n] display additional notes to words descriptions

(for regexp syntax see https://en.wikipedia.org/wiki/Regular_expression)

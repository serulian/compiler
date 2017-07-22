$module('templatestr', function () {
  var $static = this;
  $static.TEST = function () {
    var a;
    var b;
    var result;
    a = $t.fastbox(1, $g.________testlib.basictypes.Integer);
    b = $t.fastbox(true, $g.________testlib.basictypes.Boolean);
    result = $g.________testlib.basictypes.formatTemplateString($g.________testlib.basictypes.Slice($g.________testlib.basictypes.String).overArray([$t.fastbox("This function is #", $g.________testlib.basictypes.String), $t.fastbox("! It is ", $g.________testlib.basictypes.String), $t.fastbox("!", $g.________testlib.basictypes.String)]), $g.________testlib.basictypes.Slice($g.________testlib.basictypes.Stringable).overArray([a, b]));
    return $g.________testlib.basictypes.String.$equals(result, $t.fastbox('This function is #1! It is true!', $g.________testlib.basictypes.String));
  };
});

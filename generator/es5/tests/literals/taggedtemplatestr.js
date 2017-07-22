$module('taggedtemplatestr', function () {
  var $static = this;
  $static.myFunction = function (pieces, values) {
    return $t.fastbox($t.cast(values.$index($t.fastbox(0, $g.________testlib.basictypes.Integer)), $g.________testlib.basictypes.Integer, false).$wrapped + values.Length().$wrapped, $g.________testlib.basictypes.Integer);
  };
  $static.TEST = function () {
    var a;
    var b;
    var result;
    a = $t.fastbox(10, $g.________testlib.basictypes.Integer);
    b = $t.fastbox(true, $g.________testlib.basictypes.Boolean);
    result = $g.taggedtemplatestr.myFunction($g.________testlib.basictypes.Slice($g.________testlib.basictypes.String).overArray([$t.fastbox("This function is #", $g.________testlib.basictypes.String), $t.fastbox("! It is ", $g.________testlib.basictypes.String), $t.fastbox("!", $g.________testlib.basictypes.String)]), $g.________testlib.basictypes.Slice($g.________testlib.basictypes.Stringable).overArray([a, b]));
    return $t.fastbox(result.$wrapped == 12, $g.________testlib.basictypes.Boolean);
  };
});

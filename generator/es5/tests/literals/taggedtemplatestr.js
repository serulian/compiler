$module('taggedtemplatestr', function () {
  var $static = this;
  $static.myFunction = function (pieces, values) {
    return $t.fastbox($t.cast(values.$index($t.fastbox(0, $g.____testlib.basictypes.Integer)), $g.____testlib.basictypes.Integer, false).$wrapped + values.Length().$wrapped, $g.____testlib.basictypes.Boolean);
  };
  $static.TEST = function () {
    var a;
    var b;
    var result;
    a = $t.fastbox(10, $g.____testlib.basictypes.Integer);
    b = $t.fastbox(true, $g.____testlib.basictypes.Boolean);
    result = $g.taggedtemplatestr.myFunction($g.____testlib.basictypes.Slice($g.____testlib.basictypes.String).overArray([$t.fastbox("This function is #", $g.____testlib.basictypes.String), $t.fastbox("! It is ", $g.____testlib.basictypes.String), $t.fastbox("!", $g.____testlib.basictypes.String)]), $g.____testlib.basictypes.Slice($g.____testlib.basictypes.Stringable).overArray([a, b]));
    return $t.fastbox(result.$wrapped == 12, $g.____testlib.basictypes.Boolean);
  };
});

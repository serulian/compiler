$module('templatestr', function () {
  var $static = this;
  $static.TEST = function () {
    var a;
    var b;
    var result;
    a = $t.fastbox(1, $g.____testlib.basictypes.Integer);
    b = $t.fastbox(true, $g.____testlib.basictypes.Boolean);
    result = $g.____testlib.basictypes.formatTemplateString($g.____testlib.basictypes.Slice($g.____testlib.basictypes.String).overArray([$t.fastbox("This function is #", $g.____testlib.basictypes.String), $t.fastbox("! It is ", $g.____testlib.basictypes.String), $t.fastbox("!", $g.____testlib.basictypes.String)]), $g.____testlib.basictypes.Slice($g.____testlib.basictypes.Stringable).overArray([a, b]));
    return $g.____testlib.basictypes.String.$equals(result, $t.fastbox('This function is #1! It is true!', $g.____testlib.basictypes.String));
  };
});

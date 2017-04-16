$module('list', function () {
  var $static = this;
  $static.TEST = function () {
    var l;
    l = $g.____testlib.basictypes.Slice($t.struct).overArray([$t.fastbox(1, $g.____testlib.basictypes.Integer), $t.fastbox(2, $g.____testlib.basictypes.Integer), $t.fastbox(3, $g.____testlib.basictypes.Integer), $t.fastbox(true, $g.____testlib.basictypes.Boolean)]);
    return $t.fastbox(l.Length().$wrapped == 4, $g.____testlib.basictypes.Boolean);
  };
});

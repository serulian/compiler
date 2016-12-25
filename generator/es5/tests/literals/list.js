$module('list', function () {
  var $static = this;
  $static.TEST = function () {
    var l;
    l = $g.____testlib.basictypes.List($t.struct).forArray([$t.fastbox(1, $g.____testlib.basictypes.Integer), $t.fastbox(2, $g.____testlib.basictypes.Integer), $t.fastbox(3, $g.____testlib.basictypes.Integer), $t.fastbox(true, $g.____testlib.basictypes.Boolean)]);
    return $t.fastbox(l.Count().$wrapped == 4, $g.____testlib.basictypes.Boolean);
  };
});

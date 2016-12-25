$module('sliceexpr', function () {
  var $static = this;
  $static.TEST = function () {
    return $g.____testlib.basictypes.Slice($g.____testlib.basictypes.Boolean).overArray([$t.fastbox(false, $g.____testlib.basictypes.Boolean), $t.fastbox(true, $g.____testlib.basictypes.Boolean), $t.fastbox(false, $g.____testlib.basictypes.Boolean)]).$index($t.fastbox(1, $g.____testlib.basictypes.Integer));
  };
});

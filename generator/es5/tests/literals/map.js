$module('map', function () {
  var $static = this;
  $static.TEST = function () {
    var map;
    map = $g.____testlib.basictypes.Map($g.____testlib.basictypes.String, $g.____testlib.basictypes.Boolean).forArrays([$t.fastbox('hello', $g.____testlib.basictypes.String), $t.fastbox('hi', $g.____testlib.basictypes.String)], [$t.fastbox(true, $g.____testlib.basictypes.Boolean), $t.fastbox(false, $g.____testlib.basictypes.Boolean)]);
    return $t.fastbox($t.syncnullcompare(map.$index($t.fastbox('hello', $g.____testlib.basictypes.String)), function () {
      return $t.fastbox(false, $g.____testlib.basictypes.Boolean);
    }).$wrapped && !$t.syncnullcompare(map.$index($t.fastbox('hi', $g.____testlib.basictypes.String)), function () {
      return $t.fastbox(true, $g.____testlib.basictypes.Boolean);
    }).$wrapped, $g.____testlib.basictypes.Boolean);
  };
});

$module('map', function () {
  var $static = this;
  $static.TEST = function () {
    var map;
    map = $g.____testlib.basictypes.Mapping($g.____testlib.basictypes.Boolean).overObject(function () {
      var obj = {
      };
      obj['hello'] = $t.fastbox(true, $g.____testlib.basictypes.Boolean);
      obj['hi'] = $t.fastbox(false, $g.____testlib.basictypes.Boolean);
      return obj;
    }());
    return $t.fastbox($t.syncnullcompare(map.$index($t.fastbox('hello', $g.____testlib.basictypes.String)), function () {
      return $t.fastbox(false, $g.____testlib.basictypes.Boolean);
    }).$wrapped && !$t.syncnullcompare(map.$index($t.fastbox('hi', $g.____testlib.basictypes.String)), function () {
      return $t.fastbox(true, $g.____testlib.basictypes.Boolean);
    }).$wrapped, $g.____testlib.basictypes.Boolean);
  };
});

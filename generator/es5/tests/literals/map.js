$module('map', function () {
  var $static = this;
  $static.TEST = function () {
    var map;
    map = $g.________testlib.basictypes.Mapping($g.________testlib.basictypes.Boolean).overObject(function () {
      var obj = {
      };
      obj['hello'] = $t.fastbox(true, $g.________testlib.basictypes.Boolean);
      obj['hi'] = $t.fastbox(false, $g.________testlib.basictypes.Boolean);
      return obj;
    }());
    return $t.fastbox($t.syncnullcompare(map.$index($t.fastbox('hello', $g.________testlib.basictypes.String)), function () {
      return $t.fastbox(false, $g.________testlib.basictypes.Boolean);
    }).$wrapped && !$t.syncnullcompare(map.$index($t.fastbox('hi', $g.________testlib.basictypes.String)), function () {
      return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
    }).$wrapped, $g.________testlib.basictypes.Boolean);
  };
});

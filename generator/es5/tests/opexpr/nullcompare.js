$module('nullcompare', function () {
  var $static = this;
  $static.TEST = function () {
    var someBool;
    someBool = null;
    return $t.syncnullcompare(someBool, function () {
      return $t.fastbox(true, $g.____testlib.basictypes.Boolean);
    });
  };
});

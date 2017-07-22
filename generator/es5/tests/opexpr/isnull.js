$module('isnull', function () {
  var $static = this;
  $static.TEST = function (a) {
    var b;
    b = $t.fastbox(1234, $g.________testlib.basictypes.Integer);
    return $t.fastbox((a == null) && (b != null), $g.________testlib.basictypes.Boolean);
  };
});

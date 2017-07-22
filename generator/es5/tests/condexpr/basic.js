$module('basic', function () {
  var $static = this;
  $static.TEST = function () {
    var value;
    value = $t.fastbox(2, $g.________testlib.basictypes.Integer);
    return value.$wrapped == 2 ? $t.fastbox(true, $g.________testlib.basictypes.Boolean) : $t.fastbox(false, $g.________testlib.basictypes.Boolean);
  };
});

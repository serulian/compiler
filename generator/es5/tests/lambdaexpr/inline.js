$module('inline', function () {
  var $static = this;
  $static.TEST = function () {
    var result;
    result = $t.fastbox(false, $g.________testlib.basictypes.Boolean);
    (function () {
      result = $t.fastbox(true, $g.________testlib.basictypes.Boolean);
      return;
    })();
    return result;
  };
});

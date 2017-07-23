$module('nativenew', function () {
  var $static = this;
  $static.TEST = function () {
    var n;
    n = $t.nativenew($global.Number)(42);
    return $t.box((n + 0) == 42, $g.________testlib.basictypes.Boolean);
  };
});

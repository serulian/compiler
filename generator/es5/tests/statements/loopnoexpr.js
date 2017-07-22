$module('loopnoexpr', function () {
  var $static = this;
  $static.TEST = function () {
    var value;
    var $current = 0;
    syncloop: while (true) {
      switch ($current) {
        case 0:
          value = $t.fastbox(false, $g.________testlib.basictypes.Boolean);
          $current = 1;
          continue syncloop;

        case 1:
          value = $t.fastbox(true, $g.________testlib.basictypes.Boolean);
          $current = 2;
          continue syncloop;

        case 2:
          return value;

        default:
          return;
      }
    }
  };
});

$module('conditional', function () {
  var $static = this;
  $static.TEST = function () {
    var $current = 0;
    syncloop: while (true) {
      switch ($current) {
        case 0:
          if (true) {
            $current = 1;
            continue syncloop;
          } else {
            $current = 2;
            continue syncloop;
          }
          break;

        case 1:
          return $t.fastbox(true, $g.________testlib.basictypes.Boolean);

        case 2:
          return $t.fastbox(false, $g.________testlib.basictypes.Boolean);

        default:
          return;
      }
    }
  };
});

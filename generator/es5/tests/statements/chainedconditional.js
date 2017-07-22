$module('chainedconditional', function () {
  var $static = this;
  $static.TEST = function () {
    var $current = 0;
    syncloop: while (true) {
      switch ($current) {
        case 0:
          if (false) {
            $current = 1;
            continue syncloop;
          } else {
            $current = 2;
            continue syncloop;
          }
          break;

        case 1:
          $t.fastbox(123, $g.________testlib.basictypes.Integer);
          return $t.fastbox(false, $g.________testlib.basictypes.Boolean);

        case 2:
          if (false) {
            $current = 3;
            continue syncloop;
          } else {
            $current = 4;
            continue syncloop;
          }
          break;

        case 3:
          $t.fastbox(456, $g.________testlib.basictypes.Integer);
          return $t.fastbox(false, $g.________testlib.basictypes.Boolean);

        case 4:
          $t.fastbox(789, $g.________testlib.basictypes.Integer);
          return $t.fastbox(true, $g.________testlib.basictypes.Boolean);

        default:
          return;
      }
    }
  };
});

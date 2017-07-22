$module('intconditional', function () {
  var $static = this;
  $static.TEST = function () {
    var first;
    var second;
    var $current = 0;
    syncloop: while (true) {
      switch ($current) {
        case 0:
          first = $t.fastbox(10, $g.________testlib.basictypes.Integer);
          second = $t.fastbox(2, $g.________testlib.basictypes.Integer);
          if (second.$wrapped <= first.$wrapped) {
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

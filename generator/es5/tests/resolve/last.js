$module('last', function () {
  var $static = this;
  $static.SomeFunction = function () {
    var a;
    var $current = 0;
    syncloop: while (true) {
      switch ($current) {
        case 0:
          try {
            var $expr = $t.fastbox(true, $g.________testlib.basictypes.Boolean);
            a = $expr;
          } catch ($rejected) {
            a = null;
          }
          $current = 1;
          continue syncloop;

        case 1:
          return;

        default:
          return;
      }
    }
  };
});

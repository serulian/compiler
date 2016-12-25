$module('switchnoexpr', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $current = 0;
    syncloop: while (true) {
      switch ($current) {
        case 0:
          $t.fastbox(123, $g.____testlib.basictypes.Integer);
          if (false) {
            $current = 1;
            continue syncloop;
          } else {
            $current = 3;
            continue syncloop;
          }
          break;

        case 1:
          $t.fastbox(1234, $g.____testlib.basictypes.Integer);
          $current = 2;
          continue syncloop;

        case 2:
          $t.fastbox(789, $g.____testlib.basictypes.Integer);
          return;

        case 3:
          if (true) {
            $current = 4;
            continue syncloop;
          } else {
            $current = 5;
            continue syncloop;
          }
          break;

        case 4:
          $t.fastbox(2345, $g.____testlib.basictypes.Integer);
          $current = 2;
          continue syncloop;

        case 5:
          if (true) {
            $current = 6;
            continue syncloop;
          } else {
            $current = 2;
            continue syncloop;
          }
          break;

        case 6:
          $t.fastbox(3456, $g.____testlib.basictypes.Integer);
          $current = 2;
          continue syncloop;

        default:
          return;
      }
    }
  };
});

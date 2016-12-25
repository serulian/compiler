$module('loop', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $current = 0;
    syncloop: while (true) {
      switch ($current) {
        case 0:
          $t.fastbox(1234, $g.____testlib.basictypes.Integer);
          $current = 1;
          continue syncloop;

        case 1:
          $t.fastbox(1357, $g.____testlib.basictypes.Integer);
          $current = 1;
          continue syncloop;

        default:
          return;
      }
    }
  };
});

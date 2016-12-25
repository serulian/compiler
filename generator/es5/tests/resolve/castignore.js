$module('castignore', function () {
  var $static = this;
  $static.TEST = function () {
    var a;
    var somevalue;
    var $current = 0;
    syncloop: while (true) {
      switch ($current) {
        case 0:
          somevalue = $t.fastbox('hello', $g.____testlib.basictypes.String);
          try {
            var $expr = $t.cast(somevalue, $g.____testlib.basictypes.Integer, false);
            a = $expr;
          } catch ($rejected) {
            a = null;
          }
          $current = 1;
          continue syncloop;

        case 1:
          return $t.fastbox(a == null, $g.____testlib.basictypes.Boolean);

        default:
          return;
      }
    }
  };
});

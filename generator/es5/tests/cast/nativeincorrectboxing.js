$module('nativeincorrectboxing', function () {
  var $static = this;
  $static.TEST = function () {
    var err;
    var result;
    var sany;
    var $current = 0;
    syncloop: while (true) {
      switch ($current) {
        case 0:
          sany = $t.fastbox('hello world', $g.____testlib.basictypes.String);
          try {
            var $expr = $t.fastbox($t.cast(sany, $global.String, false), $g.____testlib.basictypes.String);
            result = $expr;
            err = null;
          } catch ($rejected) {
            err = $t.ensureerror($rejected);
            result = null;
          }
          $current = 1;
          continue syncloop;

        case 1:
          return $t.fastbox((result == null) && (err != null), $g.____testlib.basictypes.Boolean);

        default:
          return;
      }
    }
  };
});

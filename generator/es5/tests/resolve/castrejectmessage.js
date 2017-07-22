$module('castrejectmessage', function () {
  var $static = this;
  $static.TEST = function () {
    var b;
    var somevalue;
    var $current = 0;
    syncloop: while (true) {
      switch ($current) {
        case 0:
          somevalue = $t.fastbox('hello', $g.________testlib.basictypes.String);
          try {
            var $expr = $t.cast(somevalue, $g.________testlib.basictypes.Integer, false);
            b = null;
          } catch ($rejected) {
            b = $t.ensureerror($rejected);
          }
          $current = 1;
          continue syncloop;

        case 1:
          return $g.________testlib.basictypes.String.$equals($t.syncnullcompare($t.dynamicaccess(b, 'Message', false), function () {
            return $t.fastbox('', $g.________testlib.basictypes.String);
          }), $t.fastbox('Cannot auto-box function String() {} to function Integer() {}', $g.________testlib.basictypes.String));

        default:
          return;
      }
    }
  };
});

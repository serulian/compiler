$module('expectrejection', function () {
  var $static = this;
  this.$class('aff287b8', 'SimpleError', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.Message = $t.property(function () {
      var $this = this;
      return $t.fastbox('yo!', $g.____testlib.basictypes.String);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Message|3|268aa058": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.DoSomething = function () {
    throw $g.expectrejection.SimpleError.new();
  };
  $static.TEST = function () {
    var a;
    var b;
    var $current = 0;
    syncloop: while (true) {
      switch ($current) {
        case 0:
          try {
            var $expr = $g.expectrejection.DoSomething();
            a = $expr;
            b = null;
          } catch ($rejected) {
            b = $t.ensureerror($rejected);
            a = null;
          }
          $current = 1;
          continue syncloop;

        case 1:
          return $t.fastbox((a == null) && $g.____testlib.basictypes.String.$equals($t.assertnotnull(b).Message(), $t.fastbox('yo!', $g.____testlib.basictypes.String)).$wrapped, $g.____testlib.basictypes.Boolean);

        default:
          return;
      }
    }
  };
});

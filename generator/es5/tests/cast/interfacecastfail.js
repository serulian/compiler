$module('interfacecastfail', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return $promise.resolve(instance);
    };
    $instance.SomeValue = $t.property(function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($t.box(2, $g.____testlib.basictypes.Integer));
        return;
      };
      return $promise.new($continue);
    });
    this.$typesig = function () {
      var computed = $t.createtypesig(['SomeValue', 3, $g.____testlib.basictypes.Integer.$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.interfacecastfail.SomeClass).$typeref()]);
      this.$typesig = function () {
        return computed;
      };
      return computed;
    };
  });

  this.$interface('SomeInterface', false, '', function () {
    var $static = this;
    this.$typesig = function () {
      var computed = $t.createtypesig(['SomeValue', 3, $g.____testlib.basictypes.Boolean.$typeref()]);
      this.$typesig = function () {
        return computed;
      };
      return computed;
    };
  });

  $static.TEST = function () {
    var $result;
    var sc;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.interfacecastfail.SomeClass.new().then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            sc = $result;
            $t.cast(sc, $g.interfacecastfail.SomeInterface, false);
            $resolve();
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});

$module('structnew', function () {
  var $static = this;
  this.$class('cd38ba2a', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField) {
      var instance = new $static();
      instance.SomeField = SomeField;
      instance.anotherField = $t.fastbox(false, $g.____testlib.basictypes.Boolean);
      return $promise.resolve(instance);
    };
    $instance.set$AnotherField = function (val) {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $this.anotherField = val;
        $resolve();
        return;
      };
      return $promise.new($continue);
    };
    $instance.AnotherField = $t.property(function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($this.anotherField);
        return;
      };
      return $promise.new($continue);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "AnotherField|3|5ab5941e": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var $result;
    var $temp0;
    var sc;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.structnew.SomeClass.new($t.fastbox(2, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              $temp0 = $result0;
              return $temp0.set$AnotherField($t.fastbox(true, $g.____testlib.basictypes.Boolean)).then(function ($result1) {
                $result = ($temp0, $result1, $temp0);
                $current = 1;
                $continue($resolve, $reject);
                return;
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            sc = $result;
            $promise.resolve(sc.SomeField.$wrapped == 2).then(function ($result0) {
              $result = $t.fastbox($result0 && sc.anotherField.$wrapped, $g.____testlib.basictypes.Boolean);
              $current = 2;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
            $resolve($result);
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

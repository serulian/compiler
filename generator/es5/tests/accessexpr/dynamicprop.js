$module('dynamicprop', function () {
  var $static = this;
  this.$class('ab4079a0', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.value = $t.fastbox(42, $g.________testlib.basictypes.Integer);
      return instance;
    };
    $instance.set$SomeProp = function (val) {
      var $this = this;
      $this.value = val;
      return;
    };
    $instance.SomeProp = $t.property(function () {
      var $this = this;
      return $this.value;
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SomeProp|3|bb8d3aad": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = $t.markpromising(function () {
    var $result;
    var sc;
    var sca;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            sc = $g.dynamicprop.SomeClass.new();
            sc.set$SomeProp($t.fastbox(123, $g.________testlib.basictypes.Integer));
            sca = sc;
            $t.dynamicaccess(sca, 'SomeProp', true).then(function ($result1) {
              return $promise.resolve($t.cast($result1, $g.________testlib.basictypes.Integer, false).$wrapped == 123).then(function ($result0) {
                $result = $t.fastbox($result0 && (sc.SomeProp().$wrapped == 123), $g.________testlib.basictypes.Boolean);
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
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
});
